package rollup

import (
	"fmt"

	"github.com/emirpasic/gods/stacks"
	"github.com/emirpasic/gods/stacks/arraystack"
	"github.com/laizy/log"
	"github.com/laizy/web3/contract"
	"github.com/laizy/web3/utils"
	"github.com/ontology-layer-2/optimistic-rollup/binding"
)

type ClientOracle interface {
	GetPendingTxBatches() (*binding.RollupInputBatches, error)
}

//fixme: only for test now, because of the unfinsish of fetch calldata
type CalldataStack struct {
	Cache stacks.Stack
}

var DataStack = &CalldataStack{Cache: arraystack.New()}

type UploadBackend struct {
	l2client   ClientOracle
	signer     *contract.Signer
	stateChain *binding.RollupStateChain
	inputChain *binding.RollupInputChain
	quit       chan struct{}
	chainId    uint64
}

func NewUploadService(l2client ClientOracle, signer *contract.Signer, stateChain *binding.RollupStateChain, inputChain *binding.RollupInputChain, chainId uint64) *UploadBackend {
	return &UploadBackend{l2client, signer, stateChain, inputChain, make(chan struct{}), chainId}
}

func (self *UploadBackend) uploadState() error {
	inputChainHeight, err := self.inputChain.ChainHeight()
	if err != nil {
		return err
	}
	stateChainHeight, err := self.stateChain.TotalSubmittedState()
	if err != nil {
		return err
	}
	if stateChainHeight < inputChainHeight {
		blocks := make([][32]byte, inputChainHeight-stateChainHeight)
		receipt := self.stateChain.AppendStateBatch(blocks, stateChainHeight).Sign(self.signer).SendTransaction(self.signer)
		if receipt.IsReverted() {
			return fmt.Errorf("append state batch failed: %s", utils.JsonString(receipt))
		}
	}
	return nil
}

func (self *UploadBackend) Start() error {
	go self.runTxTask()
	return nil
}

//fixme: now only support one sequencer.
func (self *UploadBackend) runTxTask() {
	for {
		select {
		case <-self.quit:
			return
		default:
		}
		if batches, err := self.l2client.GetPendingTxBatches(); err != nil {
			log.Error(err.Error())
			continue
		} else {
			receipt := self.inputChain.AppendInputBatches(batches).Sign(self.signer).SendTransaction(self.signer)
			if receipt.IsReverted() {
				log.Errorf("append input batch failed: %s", utils.JsonString(receipt))
			}
		}
	}
}

func (self *UploadBackend) Stop() error {
	self.quit <- struct{}{}
	return nil
}
