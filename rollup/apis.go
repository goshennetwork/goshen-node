package rollup

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ontology-layer-2/optimistic-rollup/binding"
)

type L2Api struct {
	*RollupBackend
}

func Apis(backend *RollupBackend) []rpc.API {
	return []rpc.API{
		{
			Namespace: "l2",
			Version:   "1.0",
			Service:   newUploadApi(backend),
			Public:    true,
		},
	}
}

func newUploadApi(rollupBackend *RollupBackend) *L2Api {
	return &L2Api{rollupBackend}
}

//fixme: now only support one sequencer.
func (self *L2Api) GetPendingTxBatches() (*binding.RollupInputBatches, error) {
	if !self.IsSynced() {
		return nil, fmt.Errorf("syncing")
	}
	l2Store := self.store.L2Client()
	info := self.store.InputChain().GetInfo()
	checkedBatchNum := l2Store.GetTotalCheckedBatchNum()
	if checkedBatchNum < info.TotalBatches {
		return nil, fmt.Errorf("nothing to append.maybe waitting to check")
	}
	l2BlockNum := l2Store.GetTotalCheckedBlockNum(info.TotalBatches)
	localNumber := self.ethBackend.BlockChain().CurrentHeader().Number.Uint64()
	if l2BlockNum > localNumber {
		//local have nothing to upload
		return nil, fmt.Errorf("nothing need to upload total uploaded block: %d, local block number: %d", l2BlockNum, localNumber)
	}
	floor := localNumber - l2BlockNum + 1
	//todo: now simple limit upload size.should limit calldata size instead
	if floor > 512 {
		floor = 512
	}
	batches := &binding.RollupInputBatches{QueueStart: info.PendingQueueIndex}
	for i := uint64(0); i < floor; i++ {
		blockNumber := i + l2BlockNum
		block := self.ethBackend.BlockChain().GetBlockByNumber(blockNumber)
		if block == nil {
			return nil, fmt.Errorf("no block exist: %d", blockNumber)
		}
		txs := block.Transactions()
		l2txs, queueNum := FilterOutQueues(txs)
		batches.QueueNum += queueNum
		if len(l2txs) > 0 {
			batches.SubBatches = append(batches.SubBatches, &binding.SubBatch{block.Time(), l2txs})
		}
	}
	return batches, nil
}

func FilterOutQueues(txs []*types.Transaction) ([]*types.Transaction, uint64) {
	ret := make([]*types.Transaction, 0, len(txs))
	queueNum := uint64(0)
	for _, tx := range txs {
		if tx.IsQueue() {
			queueNum++
		} else {
			ret = append(ret, tx)
		}
	}
	return ret, queueNum
}
