package ethapi

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/consts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rollup"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ontology-layer-2/rollup-contracts/binding"
	"github.com/ontology-layer-2/rollup-contracts/store/schema"
)

type L2Api struct {
	*rollup.RollupBackend
}

func Apis(backend *rollup.RollupBackend) []rpc.API {
	return []rpc.API{
		{
			Namespace: "l2",
			Version:   "1.0",
			Service:   newL2Api(backend),
			Public:    true,
		},
	}
}

func newL2Api(rollupBackend *rollup.RollupBackend) *L2Api {
	return &L2Api{rollupBackend}
}

type GlobalInfo struct {
	//total batch num in l1 RollupInputChain contract
	L1InputInfo schema.InputChainInfo
	//l2 client have checked tx batch num
	L2CheckedBatchNum uint64
	//the total block num l2 already checked,start from 1, because genesis block do not need to check
	L2CheckedBlockNum uint64
	//l2 client head block num
	L2HeadBlockNumber   uint64
	L1SyncedBlockNumber uint64
	L1SyncedTimestamp   *uint64
}

func (self *L2Api) GlobalInfo() *GlobalInfo {
	//maybe syncing should also return some info
	var ret GlobalInfo
	l2Store := self.Store.L2Client()
	info := self.Store.InputChain().GetInfo()
	ret.L1InputInfo = *info
	ret.L2CheckedBatchNum = l2Store.GetTotalCheckedBatchNum()
	ret.L2CheckedBlockNum = l2Store.GetTotalCheckedBlockNum(ret.L2CheckedBatchNum - 1)
	ret.L2HeadBlockNumber = self.EthBackend.BlockChain().CurrentHeader().Number.Uint64()
	ret.L1SyncedBlockNumber = self.Store.GetLastSyncedL1Height()
	ret.L1SyncedTimestamp = self.Store.GetLastSyncedL1Timestamp()
	return &ret
}

//fixme: now only support one sequencer.
// GetPendingTxBatches return the batchCode, which is params in AppendBatch, set func selector in front of it to invoke
//append input batch.
func (self *L2Api) GetPendingTxBatches() []byte {
	if !self.IsSynced() {
		log.Warn("syncing")
		return nil
	}
	info := self.GlobalInfo()
	if info.L2CheckedBatchNum < info.L1InputInfo.TotalBatches { //should check all the batch first
		log.Warn("nothing to append")
		return nil
	}
	l2CheckedBlockNum, l2HeadBlockNumber := info.L2CheckedBlockNum, info.L2HeadBlockNumber
	if l2CheckedBlockNum > l2HeadBlockNumber {
		//local have nothing to upload
		log.Warn("nothing need to upload ", "total checked block", l2CheckedBlockNum, "local block number", l2HeadBlockNumber)
		return nil
	}
	maxBlockes := l2HeadBlockNumber - l2CheckedBlockNum + 1
	//todo: now simple limit upload size.should limit calldata size instead
	if maxBlockes > 512 {
		maxBlockes = 512
	}
	batches := &binding.RollupInputBatches{QueueStart: info.L1InputInfo.PendingQueueIndex, BatchIndex: info.L2CheckedBatchNum}
	var batchesData []byte
	for i := uint64(0); i < maxBlockes; i++ {
		blockNumber := i + l2CheckedBlockNum
		block := self.EthBackend.BlockChain().GetBlockByNumber(blockNumber)
		if block == nil { // should not happen except chain reorg
			log.Warn("nil block", "blockNumber", blockNumber)
			return nil
		}
		txs := block.Transactions()
		l2txs, queueNum := FilterOutQueues(txs)
		batches.QueueNum += queueNum
		if len(l2txs) > 0 {
			batches.SubBatches = append(batches.SubBatches, &binding.SubBatch{Timestamp: block.Time(), Txs: l2txs})
		}
		newBatch := batches.Encode()
		if len(newBatch)+4 < consts.MaxRollupInputBatchSize {
			batchesData = newBatch
		}
	}
	log.Info("generate batch", "index", batches.BatchIndex, "size", len(batchesData))
	return batchesData
}

func (self *L2Api) GetState(batchIndex uint64) common.Hash {
	if batchIndex >= self.Store.L2Client().GetTotalCheckedBatchNum() {
		return common.Hash{}
	}
	blockNum := self.Store.L2Client().GetTotalCheckedBlockNum(batchIndex)
	index := blockNum - 1
	block := self.EthBackend.BlockChain().GetBlockByNumber(index)
	if block == nil {
		log.Warn("nil block", "blockNumber", index)
		return common.Hash{}
	}
	return block.Hash()
}

// InputBatchNumber return the latest input batch number of L1
func (self *L2Api) InputBatchNumber() (uint64, error) {
	info, err := self.RollupBackend.LatestInputBatchInfo()
	if err != nil {
		return 0, err
	}
	return info.TotalBatches, nil
}

// StateBatchNumber return the latest input batch number of L1
func (self *L2Api) StateBatchNumber() (uint64, error) {
	info, err := self.RollupBackend.LatestStateBatchInfo()
	if err != nil {
		return 0, err
	}
	return info.TotalSize, nil
}

type RPCEnqueuedTx struct {
	QueueIndex hexutil.Uint64 `json:"queueIndex"`
	From       common.Address `json:"from"`
	To         common.Address `json:"to"`
	RlpTx      hexutil.Bytes  `json:"rlpTx"`
	Timestamp  hexutil.Uint64 `json:"timestamp"`
}

func (self *L2Api) GetEnqueuedTxs(queueStart, queueNum uint64) ([]*RPCEnqueuedTx, error) {
	txs, err := self.RollupBackend.GetEnqueuedTxs(queueStart, queueNum)
	if err != nil {
		return nil, err
	}
	results := make([]*RPCEnqueuedTx, 0)
	for _, tx := range txs {
		results = append(results, &RPCEnqueuedTx{
			QueueIndex: hexutil.Uint64(tx.QueueIndex),
			From:       common.Address(tx.From),
			To:         common.Address(tx.To),
			RlpTx:      tx.RlpTx,
			Timestamp:  hexutil.Uint64(tx.Timestamp),
		})
	}
	return results, nil
}

// GetBatch return the detail of batch input
func (self *L2Api) GetBatch(batchNumber uint64, useDetail bool) (map[string]interface{}, error) {
	batch, err := self.RollupBackend.InputBatchByNumber(batchNumber)
	if err != nil {
		return nil, err
	}
	batchData, err := self.RollupBackend.InputBatchDataByNumber(batchNumber)
	if err != nil {
		return nil, err
	}
	result := make(map[string]interface{}, 0)
	result["sequencer"] = batch.Proposer.String()
	result["batchNumber"] = hexutil.Uint64(batch.Index)
	result["batchHash"] = batch.InputHash.String()
	result["queueNum"] = hexutil.Uint64(batch.QueueNum)
	result["queueStart"] = hexutil.Uint64(batch.StartQueueIndex)
	formatTx := func(tx *types.Transaction) interface{} {
		return tx.Hash()
	}
	if useDetail {
		formatTx = func(tx *types.Transaction) interface{} {
			tx, blockHash, blockNumber, index := rawdb.ReadTransaction(self.EthBackend.ChainDb(), tx.Hash())
			// base fee is nil
			return newRPCTransaction(tx, blockHash, blockNumber, index, nil, self.EthBackend.BlockChain().Config())
		}
	}
	transactions := make([]interface{}, 0)
	for _, subBatchs := range batchData.SubBatches {
		for _, tx := range subBatchs.Txs {
			formatedTx := formatTx(tx)
			transactions = append(transactions, formatedTx)
		}
	}
	result["transactions"] = transactions
	return result, nil
}

type RPCBatchState struct {
	Index     hexutil.Uint64
	Proposer  common.Address
	Timestamp hexutil.Uint64
	BlockHash common.Hash
}

// GetBatchState return the state of batch input
func (self *L2Api) GetBatchState(batchNumber uint64) (*RPCBatchState, error) {
	batchState, err := self.RollupBackend.BatchState(batchNumber)
	if err != nil {
		return nil, err
	}
	return &RPCBatchState{
		Index:     hexutil.Uint64(batchState.Index),
		Proposer:  common.Address(batchState.Proposer),
		Timestamp: hexutil.Uint64(batchState.Timestamp),
		BlockHash: common.Hash(batchState.BlockHash),
	}, nil
}

func (self *L2Api) GetL2MMRProof(msgIndex, size uint64) ([]common.Hash, error) {
	proof, err := self.RollupBackend.GetL2MMRProof(msgIndex, size)
	if err != nil {
		return nil, err
	}
	result := make([]common.Hash, len(proof))
	for i := range proof {
		result[i] = common.Hash(proof[i])
	}
	return result, nil
}

type L1RelayMsgParams struct {
	Target       common.Address `json:"target"`
	Sender       common.Address `json:"sender"`
	Message      hexutil.Bytes  `json:"message"`
	MessageIndex uint64         `json:"messageIndex"`
	RLPHeader    hexutil.Bytes  `json:"rlpHeader"`
	StateInfo    *RPCBatchState `json:"stateInfo"`
	Proof        []common.Hash  `json:"proof"`
}

func (self *L2Api) GetL1RelayMsgParams(msgIndex uint64) (*L1RelayMsgParams, error) {
	result := &L1RelayMsgParams{MessageIndex: msgIndex}
	msg, err := self.RollupBackend.GetL2SentMessage(msgIndex)
	if err != nil {
		return nil, err
	}
	result.Target = common.Address(msg.Target)
	result.Sender = common.Address(msg.Sender)
	result.Message = msg.Message
	blockNum := msg.BlockNumber + 1
	batchNum := self.RollupBackend.GetL2BlockNumToBatchNum(blockNum)
	stateInfo, err := self.GetBatchState(batchNum)
	if err != nil {
		return nil, err
	}
	result.StateInfo = stateInfo
	header := self.EthBackend.BlockChain().GetHeaderByHash(stateInfo.BlockHash)
	rlpHeader, err := rlp.EncodeToBytes(header)
	if err != nil {
		return nil, fmt.Errorf("encode header, %s", err)
	}
	result.RLPHeader = rlpHeader
	// header.nonce is MMR size
	proofs, err := self.RollupBackend.GetL2MMRProof(msgIndex, header.Nonce.Uint64())
	if err != nil {
		return nil, err
	}
	result.Proof = make([]common.Hash, len(proofs))
	for i := range proofs {
		result.Proof[i] = common.Hash(proofs[i])
	}
	return result, nil
}

type L2RelayMsgParams struct {
	Target       common.Address `json:"target"`
	Sender       common.Address `json:"sender"`
	Message      hexutil.Bytes  `json:"message"`
	MessageIndex uint64         `json:"messageIndex"`
	MMRSize      uint64         `json:"mmrSize"`
	Proof        []common.Hash  `json:"proof"`
}

func (self *L2Api) GetL2RelayMsgParams(msgIndex uint64) (*L2RelayMsgParams, error) {
	result := &L2RelayMsgParams{MessageIndex: msgIndex}
	msg, err := self.RollupBackend.GetL1SentMessage(msgIndex)
	if err != nil {
		return nil, err
	}
	result.Target = common.Address(msg.Target)
	result.Sender = common.Address(msg.Sender)
	result.Message = msg.Message
	// index of l1 sent msg is MMR size
	result.MMRSize = msgIndex + 1
	proofs, err := self.RollupBackend.GetL1MMRProof(msgIndex, result.MMRSize)
	if err != nil {
		return nil, err
	}
	result.Proof = make([]common.Hash, len(proofs))
	for i := range proofs {
		result.Proof[i] = common.Hash(proofs[i])
	}
	return result, nil
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
