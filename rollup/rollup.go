package rollup

import (
	"errors"
	"fmt"
	"github.com/ontology-layer-2/rollup-contracts/blob"
	"sort"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/goshennetwork/rollup-contracts/binding"
	"github.com/goshennetwork/rollup-contracts/store"
	"github.com/goshennetwork/rollup-contracts/store/schema"
	lru "github.com/hashicorp/golang-lru"
	"github.com/laizy/web3"
	"github.com/laizy/web3/jsonrpc"
)

type TxsWithContext struct {
	Txs       []*types.Transaction
	Timestamp uint64
}

type EthBackend interface {
	BlockChain() *core.BlockChain
	ChainDb() ethdb.Database
	TxPool() *core.TxPool
}

type RollupBackend struct {
	EthBackend EthBackend
	Store      *store.Storage
	blobOracle blob.BlobOracle
	//l1 client
	L1Client   *jsonrpc.Client
	IsVerifier bool
	blockCache *lru.Cache
}

func NewBackend(ethBackend EthBackend, db schema.PersistStore, l1client *jsonrpc.Client, isVerifier bool) *RollupBackend {
	cache, _ := lru.New(1024)
	return &RollupBackend{ethBackend, store.NewStorage(db), nil, l1client, isVerifier, cache}
}

func (self *RollupBackend) IsSynced() bool {
	syncedHeight := self.Store.GetLastSyncedL1Height()
	l1Height, err := self.L1Client.Eth().BlockNumber()
	if err != nil {
		//network tolerate
		log.Warn("query l1 client failed", "err", err)
		return false
	}
	//more flex
	return syncedHeight+32 >= l1Height
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
	return self.Store.L1MMR().GetCompactMerkleTree().InclusionProof(msgIndex, size)
}

func (self *RollupBackend) GetL2BlockNumToBatchIndex(blockNumber uint64) (uint64, error) {
	blockNum := blockNumber + 1
	totalBatch := self.Store.L2Client().GetTotalCheckedBatchNum()
	checkedBlockNum := self.Store.L2Client().GetTotalCheckedBlockNum(totalBatch - 1)
	if blockNum > checkedBlockNum {
		return 0, fmt.Errorf("have not checked yet: blockNumber: %d, checkedBlockHeight: %d", blockNum-1, checkedBlockNum-1)
	}
	if v, exist := self.blockCache.Get(blockNum); exist {
		return v.(uint64), nil
	}
	if uint64(int(totalBatch)) != totalBatch || int(totalBatch) < 0 { //check type conversion safe
		return 0, fmt.Errorf("totalBatch too large: totalBatch: %d", totalBatch)
	}
	index := sort.Search(int(totalBatch), func(i int) bool {
		return self.Store.L2Client().GetTotalCheckedBlockNum(uint64(i)) >= blockNum
	})
	self.blockCache.Add(blockNum, uint64(index))
	return uint64(index), nil
}

func (self *RollupBackend) GetL2SentMessage(msgIndex uint64) (*schema.CrossLayerSentMessage, error) {
	return self.Store.L2CrossLayerWitness().GetSentMessage(msgIndex)
}

func (self *RollupBackend) GetL2MMRProof(msgIndex, size uint64) ([]web3.Hash, error) {
	return self.Store.L2MMR().GetCompactMerkleTree().InclusionProof(msgIndex, size)
}
