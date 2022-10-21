package rollup

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/laizy/log"
	"github.com/laizy/web3/utils"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/store/schema"
)

// CachedHeaderChain for l2 consensus engine prepare, only GetHeader wanted
type CachedHeaderChain []*types.Header

func newCachedHeaderChain(cap int) *CachedHeaderChain {
	var c CachedHeaderChain = make([]*types.Header, 0, cap)
	return &c
}
func (c *CachedHeaderChain) append(h *types.Header) {
	*c = append(*c, h)
}

func (c *CachedHeaderChain) Config() *params.ChainConfig {
	panic(1)
}
func (c *CachedHeaderChain) CurrentHeader() *types.Header {
	panic(1)
}

func (c *CachedHeaderChain) GetHeader(hash common.Hash, number uint64) *types.Header {
	for _, h := range *c {
		if h.Hash() == hash && h.Number.Uint64() == number {
			return h
		}
	}
	return nil
}

func (c *CachedHeaderChain) GetHeaderByNumber(number uint64) *types.Header {
	panic(1)
}

func (c *CachedHeaderChain) GetHeaderByHash(hash common.Hash) *types.Header {
	panic(1)
}

type WitnessService struct {
	*RollupBackend
	quit chan struct{}
}

func NewWitnessService(backend *RollupBackend) *WitnessService {
	return &WitnessService{backend, make(chan struct{})}
}

func (self *WitnessService) Start() error {
	go self.run()
	return nil
}

func (self *WitnessService) run() {
	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-ticker.C:
			if err := self.Work(); err != nil {
				log.Error("check batch failed", "err", err)
			}
		case <-self.quit:
			return
		}
	}
}

func (self *WitnessService) Save(blocks []*BlockWithReceipts) {
	chain := self.EthBackend.BlockChain()
	for _, blockWithReceipts := range blocks {
		block := blockWithReceipts.b
		savedB := chain.GetBlockByNumber(block.NumberU64())
		if savedB != nil {
			if savedB.Hash() != block.Hash() {
				log.Warn("witness found malicious block", "blockNumber", savedB.Number(), "blockHash", savedB.Hash())
			} else { //same block do not try to writeBlockWithState, which will reset head block to this block
				continue
			}
		}
		var logs []*types.Log
		for i, receipt := range blockWithReceipts.r {
			// add block location fields
			receipt.BlockHash = block.Hash()
			receipt.BlockNumber = block.Number()
			receipt.TransactionIndex = uint(i)

			// Update the block hash in all logs since it is now available and not when the
			// receipt/log of individual transactions were created.
			for _, l := range receipt.Logs {
				l.BlockHash = block.Hash()
			}
			logs = append(logs, receipt.Logs...)
		}
		// Commit block and state to database.
		_, err := self.EthBackend.BlockChain().WriteBlockWithState(block, blockWithReceipts.r, logs, blockWithReceipts.s, true, true)
		if err != nil { // wired
			panic(err)
		}

	}
}

func (self *WitnessService) Work() error {
	highestL1CheckPoint := self.Store.GetHighestL1CheckPointInfo()
	if highestL1CheckPoint == nil {
		return errors.New("no l1 checkpoint yet")
	}
	startPoint := highestL1CheckPoint.StartPoint
	if startPoint < 1 {
		return errors.New("l1 checkpoint not old enough")
	}

	localVersion := self.Store.L2Client().GetVersion()
	dbVersion := self.Store.GetL1DbVersion()
	rollbackBlock := false
	// something happend when reuse an old rollupdb, when setup, local db version is zero, which may differ from old, just return 0, to reorg
	info := self.Store.L2Client().GetPendingCheckPoint()
	if highestL1CheckPoint.StartPoint < info.StartPoint { // should never happen
		panic(1)
	}
	var oldTxs []*types.Transaction
	if dbVersion != localVersion {
		if dbVersion < localVersion { //should never happen
			panic(1)
		}
		log.Info("different version found, try to rollback node state...")
		// handle batch
		//every version gap just ez return to last version until version 0
		// i.e. db version 0 -> 1 -> 2 -> 3 -> 4 -> 5   l2client version last index 4 -> db version 1;now version up to 5 means
		// the worst db situation is that db state actual roll back to version 0.
		rollbackBlock = true
		//now just write totalBatches to specific batchIndex
		writer := self.Store.Writer()
		writer.L2Client().StoreTotalCheckedBatchNum(info.BatchIndex + 1 - 1) // start from last batch
		writer.Commit()

		//handle l2 blocks
	}
	sequencerRollBack := rollbackBlock && !self.IsVerifier

	// after roll back, try to move old tx when roll back finished, only self is miner
	if sequencerRollBack {
		head := self.EthBackend.BlockChain().CurrentHeader()
		checkedBatchNum := self.Store.L2Client().GetTotalCheckedBatchNum()
		lastIndex := uint64(0)
		if checkedBatchNum > 0 {
			lastIndex = checkedBatchNum - 1
		}
		checkedBlockNum := self.Store.L2Client().GetTotalCheckedBlockNum(lastIndex)
		for blockNumber := checkedBlockNum; blockNumber <= head.Number.Uint64(); blockNumber++ {
			oldTxs = append(oldTxs, self.EthBackend.BlockChain().GetBlockByNumber(blockNumber).Transactions()...)
		}
	}

	first := true
	alreadyRun := false
	for {
		writer := self.Store.Writer()
		pendingInfo := self.Store.L2Client().GetPendingCheckPoint()
		inputChainStore := self.Store.InputChain()
		l2Batches := inputChainStore.GetInfo().TotalBatches
		checkedBatchNum := self.Store.L2Client().GetTotalCheckedBatchNum()
		lastIndex := uint64(0)
		if checkedBatchNum > 0 {
			lastIndex = checkedBatchNum - 1
		}
		run := checkedBatchNum < l2Batches //there is some batch to run
		if !run && !first {                //only permit first run to go pass
			break
		}
		first = false
		if run {
			alreadyRun = true
			checkedBlockNum := self.Store.L2Client().GetTotalCheckedBlockNum(lastIndex)
			//now only check one batch
			input, err := inputChainStore.GetAppendedTransaction(checkedBatchNum)
			if err != nil {
				// not found, roll back happens
				return err
			}

			data, err := inputChainStore.GetSequencerBatchData(input.Index)
			if err != nil {
				// not found, roll back happens
				return err
			}

			parentBlock := self.EthBackend.BlockChain().GetBlockByNumber(checkedBlockNum - 1)
			//r1cs debug
			log.Warn("r1cs debug", "want block number", checkedBlockNum-1, "head block number", self.EthBackend.BlockChain().CurrentHeader().Number.Uint64(), "batchIndex", checkedBatchNum)

			blocks, inputHash, err := ProcessBatch(data, parentBlock.Hash(), self.RollupBackend)
			if err != nil { // only roll back happen occur
				return err
			}
			utils.EnsureTrue(inputHash == input.InputHash)

			if self.IsVerifier || rollbackBlock { //verifier store blocks, roll back also save blocks
				self.Save(blocks)
			} else {
				for i, blockWithReceipt := range blocks {
					local := self.EthBackend.BlockChain().GetBlockByNumber(checkedBlockNum + uint64(i))
					block := blockWithReceipt.b
					if !bytes.Equal(local.Hash().Bytes(), block.Hash().Bytes()) {
						log.Error("wrong block found", "got number", block.Number(), "local number", local.Number(), "got hash", block.Hash(), "local hash", local.Hash())
						for _, tx := range local.Transactions() {
							log.Info("local tx", "hash", tx.Hash(), "queue", tx.IsQueue())
						}
						for _, tx := range block.Transactions() {
							log.Info("seal tx", "hash", tx.Hash(), "queue", tx.IsQueue())
						}
						log.Info("header", "local", utils.JsonStr(local.Header()), "got", utils.JsonStr(block.Header()))
						return fmt.Errorf("failed")
					}
				}
			}
			checkedBlockNum += uint64(len(blocks))
			checkedBatchNum += 1
			writer.L2Client().StoreTotalCheckedBatchNum(checkedBatchNum)
			//save batch index => block num
			writer.L2Client().StoreCheckedBlockNum(checkedBatchNum-1, checkedBlockNum)
		}

		head := self.EthBackend.BlockChain().CurrentHeader()
		// every cycle read start point,to update batchindex to latest checkpoint
		startPoint = self.Store.GetHighestL1CheckPointInfo().StartPoint
		if startPoint > pendingInfo.StartPoint { //if start point update, update startPoint with batchIndex
			if run {
				log.Debug("save confirmPoint", "pending batchIndex", checkedBatchNum-1, "startpoint", startPoint, "pending block number", head.Number.Uint64())
				writer.L2Client().StorePendingCheckpoint(&schema.L2CheckPointInfo{startPoint, checkedBatchNum - 1, head.Number.Uint64()})
			} else { //move batch index to pending batchIndex, so next time will not roll back batch
				log.Debug("save confirmPoint", "pending batchIndex", checkedBatchNum, "startpoint", startPoint, "prnding block number", head.Number.Uint64())
				writer.L2Client().StorePendingCheckpoint(&schema.L2CheckPointInfo{startPoint, checkedBatchNum, head.Number.Uint64()})
			}
		}
		writer.Commit()
		self.EthBackend.TxPool().AddLocals(oldTxs)
		oldTxs = oldTxs[:0] //clean oldTxs
	}
	writer := self.Store.Writer()
	writer.L2Client().SetVersion(dbVersion) // now update local version to db Version
	writer.Commit()

	if !rollbackBlock { // no need to roll back
		return nil
	}
	if alreadyRun {
		// if have already run do not need to roll back block, because head block now is newest
		return nil
	}
	// need roll back and have not yet rolled back, try to roll back block
	return self.RollBackBlock()
}

func (self *WitnessService) RollBackBlock() error {
	pendingInfo := self.Store.L2Client().GetPendingCheckPoint()
	if pendingInfo.BlockNumber < 1 {
		return fmt.Errorf("pending info blockNumber less than 1")
	}
	totalBatch := self.Store.L2Client().GetTotalCheckedBatchNum()
	lastIndex := uint64(0)
	if totalBatch > 0 {
		lastIndex = totalBatch - 1
	}
	totalBlock := self.Store.L2Client().GetTotalCheckedBlockNum(lastIndex)
	//cut down， but ensure not cut down to checked batch number
	newHead := pendingInfo.BlockNumber - 1
	if newHead < totalBlock-1 {
		newHead = totalBlock - 1
	}
	var oldTxs []*types.Transaction
	for i := pendingInfo.BlockNumber; i <= self.EthBackend.BlockChain().CurrentHeader().Number.Uint64(); i++ {
		oldTxs = append(oldTxs, self.EthBackend.BlockChain().GetBlockByNumber(i).Transactions()...)
	}
	if err := self.EthBackend.BlockChain().SetHead(newHead); err != nil {
		return err
	}
	self.EthBackend.TxPool().AddLocals(oldTxs)
	log.Warn("cutting down chain head", "number", newHead)
	return nil
}

func ProcessBatch(input []byte, parentBlockHash common.Hash, rollupBackend *RollupBackend) ([]*BlockWithReceipts, [32]byte, error) {
	eth := rollupBackend.EthBackend
	parent := eth.BlockChain().GetBlockByHash(parentBlockHash)
	inputChainStore := rollupBackend.Store.InputChain()
	batch := &binding.RollupInputBatches{}
	if err := batch.Decode(input); err != nil { //may decode failed, just ignore
		log.Warn("malicious batch code found", "err", err)
	}
	if parent.Header().TotalExecutedQueueNum() != batch.QueueStart {
		panic(1)
	}
	queues, err := inputChainStore.GetEnqueuedTransactions(batch.QueueStart, batch.QueueNum)
	if err != nil { // may happen when roll back
		return nil, [32]byte{}, err
	}

	var orderTxs []*binding.SubBatch
	txs := &binding.SubBatch{}
	last := len(queues) - 1
	for i, enqueue := range queues {
		tx, timestamp := enqueue.MustToTransaction(), enqueue.Timestamp
		if i == 0 {
			txs.Timestamp = timestamp
			txs.Txs = append(txs.Txs, tx)
		} else {
			if timestamp != txs.Timestamp {
				//append first
				orderTxs = append(orderTxs, txs)
				txs = &binding.SubBatch{timestamp, nil}
			}
			txs.Txs = append(txs.Txs, tx)
		}
		if i == last {
			orderTxs = append(orderTxs, txs)
		}

	}
	orderTxs = append(orderTxs, batch.SubBatches...)

	//sort in timestamp
	sort.SliceStable(orderTxs, func(i, j int) bool {
		return orderTxs[i].Timestamp < orderTxs[j].Timestamp
	})
	return RunOrderesTxs(eth.BlockChain(), orderTxs, parent.Header()), batch.InputHash(schema.CalcQueueHash(queues)), nil

}

type blockTask struct {
	header   *types.Header
	statedb  *state.StateDB
	gasPool  *core.GasPool
	txs      []*types.Transaction
	receipts []*types.Receipt
}

func (b *blockTask) sealToBlock(chain *core.BlockChain) *BlockWithReceipts {
	engine := chain.Engine()
	// no uncles
	block, err := engine.FinalizeAndAssemble(chain, b.header, b.statedb, b.txs, nil, b.receipts)
	if err != nil {
		//should never happen?
		panic(err)
	}
	return &BlockWithReceipts{block, b.statedb, b.receipts}
}

// state is parent state,just use it
func makeBlockTask(parent *types.Header, queueNum uint64, timestamp uint64, stateDb *state.StateDB, chain *core.BlockChain, fakeHeaderChain *CachedHeaderChain) *blockTask {
	engine := chain.Engine()
	num := parent.Number.Uint64()
	gaslimit := parent.GasLimit
	h := &types.Header{
		ParentHash: parent.Hash(),
		Number:     new(big.Int).SetUint64(num + 1),
		GasLimit:   gaslimit,
		Extra:      nil,
		Time:       timestamp,
	}
	h.SetTotalExecutedQueueNum(parent.TotalExecutedQueueNum() + queueNum)
	err := engine.Prepare(fakeHeaderChain, h)
	utils.Ensure(err)
	return &blockTask{
		header:  h,
		statedb: stateDb.Copy(),
		gasPool: new(core.GasPool).AddGas(gaslimit),
	}
}

type BlockWithReceipts struct {
	b *types.Block
	s *state.StateDB
	r []*types.Receipt
}

func RunOrderesTxs(chain *core.BlockChain, orderTxs []*binding.SubBatch, parent *types.Header) []*BlockWithReceipts {
	var ret []*BlockWithReceipts
	fakeHeaderChain := newCachedHeaderChain(len(orderTxs) + 1)
	fakeHeaderChain.append(parent)
	statedb, err := chain.StateAt(parent.Root)
	if err != nil {
		panic(err)
	}
	cfg := core.ValidDataTxConfig{
		MaxGas:   parent.GasLimit,
		GasPrice: big.NewInt(0),
		BaseFee:  nil,
		Istanbul: true,
		Eip2718:  true,
		Eip1559:  false,
	}

	gasLimit := parent.GasLimit
	parentStateDb := statedb
	for i := 0; i < len(orderTxs); i++ {
		timestamp := orderTxs[i].Timestamp
		isQueue := false
		if len(orderTxs[i].Txs) > 0 {
			isQueue = orderTxs[i].Txs[0].IsQueue()
		}
		usedGas := uint64(0)
		var txs []*types.Transaction
		var queues []*types.Transaction
		for _, tx := range orderTxs[i].Txs {
			switch isQueue {
			case true: //queue tx
				//no matter what tx consumed,just assume it consumes all gas
				usedGas += tx.Gas()
				if usedGas >= gasLimit { // a new block tak
					blockTask := makeBlockTask(parent, uint64(len(queues)), timestamp, parentStateDb, chain, fakeHeaderChain)
					CommitTransactions(chain, queues, blockTask)
					// save blocks which have executed queues no matter what happened
					ret = append(ret, blockTask.sealToBlock(chain))
					parentB := ret[len(ret)-1]
					parent = parentB.b.Header()
					fakeHeaderChain.append(parent)
					parentStateDb = blockTask.statedb
					if tx.Gas() >= gasLimit { //should never happen
						panic(1)
					}
					usedGas = tx.Gas()
					queues = queues[:0]
				}
				queues = append(queues, tx)
			case false: //l2 origin tx
				//check l2 origin tx first
				if err := core.ValidateTx(tx, parentStateDb, types.LatestSigner(chain.Config()), cfg, true); err != nil {
					log.Error("validate tx failed", "hash", tx.Hash(), "err", err, "txs", txs)
					continue
				}
				txs = append(txs, tx)
			}
		}
		blockTask := makeBlockTask(parent, uint64(len(queues)), timestamp, parentStateDb, chain, fakeHeaderChain)
		txs = append(queues, txs...)
		CommitTransactions(chain, txs, blockTask)
		//every order txs try to seal to a block
		ret = append(ret, blockTask.sealToBlock(chain))
		parentB := ret[len(ret)-1]
		parent = parentB.b.Header()
		fakeHeaderChain.append(parent)
		parentStateDb = blockTask.statedb
	}
	return ret
}

func (b *blockTask) TotalNonQueueTxSize() uint64 {
	s := uint64(0)
	for _, tx := range b.txs {
		if !tx.IsQueue() {
			s += uint64(tx.Size())
		}
	}
	return s
}

func CommitTransactions(chain *core.BlockChain, txs []*types.Transaction, task *blockTask) {
	log.Debug("commitTransactions", "blockNumber", task.header.Number.Uint64(), "txNum", len(txs))
	txIndex := 0
	for _, tx := range txs {
		statedb := task.statedb
		header := task.header
		gasPool := task.gasPool
		statedb.Prepare(tx.Hash(), txIndex)
		snap := statedb.Snapshot()
		chainConfig := chain.Config()
		receipt, err := core.ApplyTransaction(chainConfig, chain, &header.Coinbase, gasPool, statedb, header, tx, &header.GasUsed, *chain.GetVMConfig())
		if err != nil {
			statedb.RevertToSnapshot(snap)
			log.Warn("malicious tx", "err", err, "txHash", tx.Hash())
			continue
		}
		task.txs = append(task.txs, tx)
		task.receipts = append(task.receipts, receipt)
		txIndex += 1
	}
}

func (self *WitnessService) Stop() error {
	self.quit <- struct{}{}
	return nil
}
