package rollup

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
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

func NewBackend(ethBackend EthBackend, db schema.PersistStore, l1client *jsonrpc.Client) *RollupBackend {
	return &RollupBackend{ethBackend, store.NewStorage(db), l1client}
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

func (self *RollupBackend) StoreAndSetExecutedNum(totalQueueChain *schema.ChainedEnqueueBlockInfo) {
	writer := self.store.Writer()
	writer.L2Client().StoreExecutedQueue(totalQueueChain.CurrEnqueueBlock, totalQueueChain)
	writer.L2Client().StoreHeadExecutedQueueBlock(totalQueueChain)
	writer.Commit()
}

func (self *RollupBackend) GetPendingQueue(parentBlock uint64, gasLimit uint64) (*QueuesWithContext,
	*schema.ChainedEnqueueBlockInfo, error) {
	if !self.IsSynced() {
		return nil, nil, errors.New("not synced yet")
	}
	l1Time := self.store.GetLastSyncedL1Timestamp()
	if l1Time == nil {
		return nil, nil, fmt.Errorf("no synced l1 timestamp")
	}
	headQueue := self.store.L2Client().GetHeadExecutedQueueBlock()
	if parentBlock < headQueue.CurrEnqueueBlock { //should never happen
		panic(1)
	}
	queuesInfo := &QueuesWithContext{}
	usedGas := uint64(0)
	pendingQueueIndex := headQueue.TotalEnqueuedTx
	for {
		enqueuedEvent, err := self.store.InputChain().GetEnqueuedTransaction(pendingQueueIndex)
		if err != nil && !errors.Is(err, schema.ErrNotFound) {
			return nil, nil, fmt.Errorf("no pending queue")
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
	headQueue.PrevEnqueueBlock = headQueue.CurrEnqueueBlock
	headQueue.CurrEnqueueBlock = parentBlock + 1
	headQueue.TotalEnqueuedTx += uint64(len(queuesInfo.Txs))
	return queuesInfo, headQueue, nil
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
