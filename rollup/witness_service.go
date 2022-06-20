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
	"github.com/ontology-layer-2/optimistic-rollup/store/schema"
	sync_service "github.com/ontology-layer-2/optimistic-rollup/sync-service"
	utils2 "github.com/ontology-layer-2/optimistic-rollup/utils"
)

type TxsSigWithContext struct {
	Sigs      []common.Hash
	Timestamp uint64
}

type WitnessService struct {
	*RollupBackend
	cfg  sync_service.Config
	quit chan struct{}
}

func NewWitnessService(backend *RollupBackend, cfg *sync_service.Config) *WitnessService {
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
	checker := &utils2.HashChecker{}
	var orderTxs []TxsSigWithContext
	txs := TxsSigWithContext{}
	for i := uint64(0); i < batch.QueueNum; i++ {
		enqueue, err := inputChainStore.GetEnqueuedTransaction(batch.QueueStart + i)
		utils.Ensure(err)
		sig, timestamp := crypto.Keccak256Hash(enqueue.RlpTx), enqueue.Timestamp
		if i == 0 {
			txs.Timestamp = timestamp
			txs.Sigs = append(txs.Sigs, sig)
		} else { //wrap different timestamp to a txsSigWithContext
			if enqueue.Timestamp != txs.Timestamp {
				orderTxs = append(orderTxs, txs)
				txs = TxsSigWithContext{[]common.Hash{sig}, timestamp}
			} else {
				txs.Sigs = append(txs.Sigs, sig)
			}
		}
	}
	if batch.QueueNum > 0 {
		//append last txSig
		orderTxs = append(orderTxs, txs)
	}
	for _, batch := range batch.SubBatches {
		txs = TxsSigWithContext{nil, batch.Timestamp}
		for _, tx := range batch.Txs {
			txs.Sigs = append(txs.Sigs, tx.Hash())
		}
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
			checker.Add(genCheckerHash(sig, timestamp, index))
			index++
		}
	}
	index = uint64(0)
	blockNumber := self.store.L2Client().GetTotalCheckedBlockNum(batchIndex)
	for i := uint64(0); i < checker.Num(); i++ {
		block := self.ethBackend.BlockChain().GetBlockByNumber(blockNumber + i)
		timestamp := block.Time()
		for _, tx := range block.Transactions() {
			checker.Meet(genCheckerHash(tx.Hash(), timestamp, index))
			index++
		}
	}
	return checker.IsEqual(), blockNumber + checker.Num()
}

func (self *WitnessService) Stop() error {
	self.quit <- struct{}{}
	return nil
}
