package rollup

import (
	"bytes"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/laizy/log"
	"github.com/laizy/web3/utils"
	"github.com/laizy/web3/utils/codec"
	"github.com/ontology-layer-2/optimistic-rollup/binding"
	"github.com/ontology-layer-2/optimistic-rollup/config"
	"github.com/ontology-layer-2/optimistic-rollup/store/schema"
)

type TxsSigWithContext struct {
	Sigs      []common.Hash
	Timestamp uint64
}

type WitnessService struct {
	*RollupBackend
	cfg  config.SyncConfig
	quit chan struct{}
}

func NewWitnessService(backend *RollupBackend, cfg *config.SyncConfig) *WitnessService {
	return &WitnessService{backend, *cfg, make(chan struct{})}
}

func (self *WitnessService) Start() error {
	go self.run()
	return nil
}

func (self *WitnessService) run() {
	ticker := time.NewTicker(1 * time.Second)
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

func (self *WitnessService) Work() error {
	inputChainStore := self.store.InputChain()
	l2Batches := inputChainStore.GetInfo().TotalBatches

	checkedBatchNum := self.store.L2Client().GetTotalCheckedBatchNum()
	if checkedBatchNum >= l2Batches {
		log.Debug("have checked all l2 batches", "checkedBatchNum", checkedBatchNum, "l2 total", l2Batches)
		return nil
	}
	//now only check one batch
	input, err := inputChainStore.GetAppendedTransaction(checkedBatchNum)
	utils.Ensure(err)

	data, err := inputChainStore.GetSequencerBatchData(input.Index)
	if err != nil {
		//not found, should never happen
		panic(1)
	}
	batch := &binding.RollupInputBatches{}
	if err := batch.Decode(data); err != nil { //Now assume only honest sequencer, should never happen
		panic(1)
	}
	queues, err := inputChainStore.GetEnqueuedTransactions(batch.QueueStart, batch.QueueNum)
	if err != nil {
		return err
	}
	batchInputHash := batch.InputHash(schema.CalcQueueHash(queues))
	utils.Ensure(err)
	//check input hash
	utils.EnsureTrue(bytes.Equal(input.InputHash.Bytes(), batchInputHash[:]))

	//check tx consistent with local, now just panic because of only honest sequencer
	ok, checkedBlockNum := self.check(batch, checkedBatchNum)
	utils.EnsureTrue(ok)
	checkedBatchNum += 1

	//now store witnessed state
	writer := self.store.Writer()
	writer.L2Client().StoreTotalCheckedBatchNum(checkedBatchNum)
	writer.L2Client().StoreCheckedBlockNum(checkedBatchNum, checkedBlockNum)
	writer.Commit()
	return nil
}

//make sure tx in order and tx timestamp is correct
func genCheckerHash(txHash [32]byte, timestamp uint64, index uint64) [32]byte {
	s := codec.NewZeroCopySink(nil)
	s.WriteHash(txHash)
	s.WriteUint64BE(timestamp)
	s.WriteUint64BE(index)
	return crypto.Keccak256Hash(s.Bytes())
}

// check batch tx with local tx, now assume sequencer is honest so txs should be all in order at local chain.
func (self *WitnessService) check(batch *binding.RollupInputBatches, batchIndex uint64) (bool, uint64) {
	inputChainStore := self.store.InputChain()
	checker := &HashChecker{}
	var orderTxs []TxsSigWithContext
	txs := TxsSigWithContext{}
	for i := uint64(0); i < batch.QueueNum; i++ {
		enqueue, err := inputChainStore.GetEnqueuedTransaction(batch.QueueStart + i)
		utils.Ensure(err)
		sig, timestamp := crypto.Keccak256Hash(enqueue.RlpTx), enqueue.Timestamp
		switch i {
		case 0:
			txs.Timestamp = timestamp
			txs.Sigs = append(txs.Sigs, sig)
		case batch.QueueNum - 1: //last
			if timestamp != txs.Timestamp {
				//append first
				orderTxs = append(orderTxs, txs)
				txs = TxsSigWithContext{nil, timestamp}
			}
			txs.Sigs = append(txs.Sigs, sig)
			//append last txSig
			orderTxs = append(orderTxs, txs)
		default:
			if timestamp != txs.Timestamp {
				orderTxs = append(orderTxs, txs)
				txs = TxsSigWithContext{nil, timestamp}
			}
			txs.Sigs = append(txs.Sigs, sig)
		}
	}

	for _, batch := range batch.SubBatches {
		txs = TxsSigWithContext{nil, batch.Timestamp}
		for _, tx := range batch.Txs {
			txs.Sigs = append(txs.Sigs, tx.Hash())
		}
		orderTxs = append(orderTxs, txs)
	}
	//sort in timestamp
	sort.SliceStable(orderTxs, func(i, j int) bool {
		return orderTxs[i].Timestamp < orderTxs[j].Timestamp
	})
	index := uint64(0)
	//now add to checker in order
	for i := 0; i < len(orderTxs); i++ {
		timestamp := orderTxs[i].Timestamp
		for _, sig := range orderTxs[i].Sigs {
			log.Debug("l2 local", "txHash", sig)
			checker.Add(genCheckerHash(sig, timestamp, index))
			index++
		}
	}
	index = uint64(0)
	blockNumber := self.store.L2Client().GetTotalCheckedBlockNum(batchIndex)
	for i := uint64(0); i < uint64(len(orderTxs)); i++ {
		block := self.ethBackend.BlockChain().GetBlockByNumber(blockNumber + i)
		timestamp := block.Time()
		for _, tx := range block.Transactions() {
			log.Debug("l2 local", "blockNumber", block.NumberU64(), "txHash", tx.Hash(), "txNonce", tx.Nonce())
			checker.Meet(genCheckerHash(tx.Hash(), timestamp, index))
			index++
		}
	}
	return checker.IsEqual(), blockNumber + uint64(len(orderTxs))
}

func (self *WitnessService) Stop() error {
	self.quit <- struct{}{}
	return nil
}

type HashChecker struct {
	task map[common.Hash]int
	num  uint64
}

//Add add filter task, duplicated hash also recorded
func (f *HashChecker) Add(h common.Hash) {
	if f.task == nil {
		f.task = make(map[common.Hash]int)
	}
	f.task[h]++
}

func (f *HashChecker) Meet(h common.Hash) {
	f.task[h]--
}

func (f *HashChecker) IsEqual() bool {
	for h, v := range f.task {
		if v != 0 {
			//debug
			log.Error("debug wrong tx", "hash", h)
			return false
		}
	}
	return true
}
