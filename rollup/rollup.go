package rollup

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/laizy/web3"
	"github.com/laizy/web3/jsonrpc"
	"github.com/ontology-layer-2/optimistic-rollup/binding"
	"github.com/ontology-layer-2/optimistic-rollup/store"
	"github.com/ontology-layer-2/optimistic-rollup/store/schema"
)

type QueuesWithContext struct {
	Txs       []*types.Transaction
	Timestamp uint64
}

type EthBackend interface {
	BlockChain() *core.BlockChain
}

type RollupBackend struct {
	ethBackend EthBackend
	store      *store.Storage
	//l1 client
	L1Client *jsonrpc.Client
}

func NewBackend(ethBackend EthBackend, db schema.PersistStore, dbPath string, l1client *jsonrpc.Client) *RollupBackend {
	return &RollupBackend{ethBackend, store.NewStorage(db, dbPath), l1client}
}

func (self *RollupBackend) IsSynced() bool {
	syncedHeight := self.store.GetLastSyncedL1Height()
	l1Height, err := self.L1Client.Eth().BlockNumber()
	if err != nil {
		//network tolerate
		log.Warn("query l1 client failed", "err", err)
		return false
	}
	//todo: avoid chain reorg event
	return syncedHeight+6 >= l1Height
}

func (self *RollupBackend) GetPendingQueue(totalExecutedQueueNum uint64, gasLimit uint64) (*QueuesWithContext, error) {
	if !self.IsSynced() {
		return nil, errors.New("not synced yet")
	}
	l1Time := self.store.GetLastSyncedL1Timestamp()
	if l1Time == nil {
		return nil, fmt.Errorf("no synced l1 timestamp")
	}
	queuesInfo := &QueuesWithContext{}
	usedGas := uint64(0)
	pendingQueueIndex := totalExecutedQueueNum
	for {
		enqueuedEvent, err := self.store.InputChain().GetEnqueuedTransaction(pendingQueueIndex)
		if err != nil && !errors.Is(err, schema.ErrNotFound) {
			return nil, fmt.Errorf("no pending queue")
		}
		if enqueuedEvent == nil {
			break
		}
		if len(queuesInfo.Txs) > 0 && enqueuedEvent.Timestamp != queuesInfo.Timestamp {
			break
		}
		tx := enqueuedEvent.MustToTransaction()
		usedGas += tx.Gas()
		if usedGas >= gasLimit {
			break
		}
		queuesInfo.Txs = append(queuesInfo.Txs, tx)
		queuesInfo.Timestamp = enqueuedEvent.Timestamp
		pendingQueueIndex++
	}
	if len(queuesInfo.Txs) == 0 {
		queuesInfo.Timestamp = *l1Time
	}

	return queuesInfo, nil
}

func (self *RollupBackend) LatestInputBatchInfo() (*schema.InputChainInfo, error) {
	return self.store.InputChain().GetInfo(), nil
}

func (self *RollupBackend) LatestStateBatchInfo() (*schema.StateChainInfo, error) {
	return self.store.StateChain().GetInfo(), nil
}

func (self *RollupBackend) InputBatchByNumber(index uint64) (*schema.AppendedTransaction, error) {
	return self.store.InputChain().GetAppendedTransaction(index)
}

func (self *RollupBackend) InputBatchDataByNumber(index uint64) (*binding.RollupInputBatches, error) {
	data, err := self.store.InputChain().GetSequencerBatchData(index)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, schema.ErrNotFound
	}
	batches := &binding.RollupInputBatches{}
	err = batches.Decode(data)
	return batches, err
}

func (self *RollupBackend) BatchState(index uint64) (*schema.RollupStateBatchInfo, error) {
	return self.store.StateChain().GetState(index)
}

func (self *RollupBackend) GetL1SentMessage(msgIndex uint64) (*schema.CrossLayerSentMessage, error) {
	return self.store.L1CrossLayerWitness().GetSentMessage(msgIndex)
}

func (self *RollupBackend) GetL1MMRProof(msgIndex, size uint64) ([]web3.Hash, error) {
	return self.store.L1CrossLayerWitness().GetL1MMRProof(msgIndex, size)
}

func (self *RollupBackend) GetL2BlockNumToBatchNum(blockNum uint64) uint64 {
	upper := self.store.L2Client().GetTotalCheckedBatchNum()
	lower := uint64(0)
	i := (upper + lower) / 2
	_block := uint64(0)
	// find the closest batch
	for {
		_block = self.store.L2Client().GetTotalCheckedBlockNum(i)
		if _block > blockNum {
			upper = i
		} else if _block < blockNum {
			lower = i
		} else {
			break
		}
		if lower >= upper-1 {
			break
		}
		i = (upper + lower) / 2
	}
	// closest higher
	if _block >= blockNum {
		return i
	}

	for {
		i++
		_block = self.store.L2Client().GetTotalCheckedBlockNum(i)
		if _block > blockNum {
			return i
		}
	}
}

func (self *RollupBackend) GetL2SentMessage(msgIndex uint64) (*schema.CrossLayerSentMessage, error) {
	return self.store.L2CrossLayerWitness().GetSentMessage(msgIndex)
}

func (self *RollupBackend) GetL2MMRProof(msgIndex, size uint64) ([]web3.Hash, error) {
	return self.store.L2CrossLayerWitness().GetL2MMRProof(msgIndex, size)
}
