package rollup

import (
	"errors"
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/laizy/web3"
	"github.com/laizy/web3/jsonrpc"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/store"
	"github.com/ontology-layer-2/rollup-contracts/store/schema"
)

type TxsWithContext struct {
	Txs       []*types.Transaction
	Timestamp uint64
}

type EthBackend interface {
	BlockChain() *core.BlockChain
	ChainDb() ethdb.Database
}

type RollupBackend struct {
	EthBackend EthBackend
	Store      *store.Storage
	//l1 client
	L1Client   *jsonrpc.Client
	IsVerifier bool
}

func NewBackend(ethBackend EthBackend, db schema.PersistStore, dbPath string, l1client *jsonrpc.Client, isVerifier bool) *RollupBackend {
	return &RollupBackend{ethBackend, store.NewStorage(db, dbPath), l1client, isVerifier}
}

func (self *RollupBackend) IsSynced() bool {
	syncedHeight := self.Store.GetLastSyncedL1Height()
	l1Height, err := self.L1Client.Eth().BlockNumber()
	if err != nil {
		//network tolerate
		log.Warn("query l1 client failed", "err", err)
		return false
	}
	//todo: avoid chain reorg event
	return syncedHeight+6 >= l1Height
}

func (self *RollupBackend) GetPendingQueue(totalExecutedQueueNum uint64, gasLimit uint64) (*TxsWithContext, error) {
	if !self.IsSynced() {
		return nil, errors.New("not synced yet")
	}
	l1Time := self.Store.GetLastSyncedL1Timestamp()
	if l1Time == nil {
		return nil, fmt.Errorf("no synced l1 timestamp")
	}
	queuesInfo := &TxsWithContext{}
	usedGas := uint64(0)
	pendingQueueIndex := totalExecutedQueueNum
	for {
		enqueuedEvent, err := self.Store.InputChain().GetEnqueuedTransaction(pendingQueueIndex)
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
	return self.Store.InputChain().GetInfo(), nil
}

func (self *RollupBackend) LatestStateBatchInfo() (*schema.StateChainInfo, error) {
	return self.Store.StateChain().GetInfo(), nil
}

func (self *RollupBackend) GetEnqueuedTxs(queueStart, queueNum uint64) ([]*schema.EnqueuedTransaction, error) {
	return self.Store.InputChain().GetEnqueuedTransactions(queueStart, queueNum)
}

func (self *RollupBackend) InputBatchByNumber(index uint64) (*schema.AppendedTransaction, error) {
	return self.Store.InputChain().GetAppendedTransaction(index)
}

func (self *RollupBackend) InputBatchDataByNumber(index uint64) (*binding.RollupInputBatches, error) {
	data, err := self.Store.InputChain().GetSequencerBatchData(index)
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
	return self.Store.StateChain().GetState(index)
}

func (self *RollupBackend) GetL1SentMessage(msgIndex uint64) (*schema.CrossLayerSentMessage, error) {
	return self.Store.L1CrossLayerWitness().GetSentMessage(msgIndex)
}

func (self *RollupBackend) GetL1MMRProof(msgIndex, size uint64) ([]web3.Hash, error) {
	return self.Store.L1CrossLayerWitness().GetL1MMRProof(msgIndex, size)
}

func (self *RollupBackend) GetL2BlockNumToBatchNum(blockNum uint64) uint64 {
	// upper is max batch index
	upper := self.Store.L2Client().GetTotalCheckedBatchNum() - 1
	index := sort.Search(int(upper), func(i int) bool {
		return self.Store.L2Client().GetTotalCheckedBlockNum(uint64(i)) >= blockNum
	})
	return uint64(index)
}

func (self *RollupBackend) GetL2SentMessage(msgIndex uint64) (*schema.CrossLayerSentMessage, error) {
	return self.Store.L2CrossLayerWitness().GetSentMessage(msgIndex)
}

func (self *RollupBackend) GetL2MMRProof(msgIndex, size uint64) ([]web3.Hash, error) {
	return self.Store.L2CrossLayerWitness().GetL2MMRProof(msgIndex, size)
}
