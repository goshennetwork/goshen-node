package rollup

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/tracers"
	"github.com/laizy/web3"
	"github.com/laizy/web3/registry"
	"github.com/laizy/web3/utils"
	"github.com/ontology-layer-2/rollup-contracts/binding"
)

func init() {
	register("rollupTracer", newLogTracer)
}

//type EVMLogger interface {
//	CaptureStart(env *EVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int)
//	CaptureState(pc uint64, op OpCode, gas, cost uint64, scope *ScopeContext, rData []byte, depth int, err error)
//	CaptureEnter(typ OpCode, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int)
//	CaptureExit(output []byte, gasUsed uint64, err error)
//	CaptureFault(pc uint64, op OpCode, gas, cost uint64, scope *ScopeContext, depth int, err error)
//	CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error)
//}

type logTracer struct {
	eventEntry *registry.EventRegistry
	env        *vm.EVM
	interrupt  uint32
	reason     error
	logs       []*types.Log
}

func newLogTracer() tracers.Tracer {
	l2Addrs := map[web3.Address]string{
		web3.HexToAddress("0x2210000000000000000000000000000000000221"): "L2CrossLayerWitness",
		web3.HexToAddress("0xbde0000000000000000000000000000000000bde"): "L2StandardBridge",
		web3.HexToAddress("0xfee0000000000000000000000000000000000fee"): "L2FeeCollector",
	}
	entry := registry.NewEventRegistry()
	entry.RegisterFromAbi(binding.L2StandardBridgeAbi())
	entry.RegisterFromAbi(binding.L2CrossLayerWitnessAbi())
	entry.RegisterFromAbi(binding.L2FeeCollectorAbi())
	entry.RegisterFromAbi(binding.ERC20Abi())
	for addr, alia := range l2Addrs {
		entry.RegisterContractAlias(addr, alia)
	}
	return &logTracer{entry, nil, 0, nil, nil}
}

func (tracer *logTracer) CaptureStart(env *vm.EVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) {
	tracer.env = env
}
func (tracer *logTracer) CaptureState(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, rData []byte, depth int, err error) {
	if atomic.LoadUint32(&tracer.interrupt) == 1 {
		tracer.env.Cancel()
		return
	}
	size := 0
	switch op {
	case vm.LOG0:
		size = 0
	case vm.LOG1:
		size = 1
	case vm.LOG2:
		size = 2
	case vm.LOG3:
		size = 3
	case vm.LOG4:
		size = 4
	default:
		return
	}
	tracer.logs = append(tracer.logs, makeLog(size)(pc, scope))
}
func (tracer *logTracer) CaptureEnter(typ vm.OpCode, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int) {
}
func (tracer *logTracer) CaptureExit(output []byte, gasUsed uint64, err error) {}
func (tracer *logTracer) CaptureFault(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, depth int, err error) {
}
func (tracer *logTracer) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) {}
func (tracer *logTracer) GetResult() (json.RawMessage, error) {
	var ret []*web3.ParsedEvent
	for _, l := range tracer.logs {
		topics := make([]web3.Hash, len(l.Topics))
		for i, v := range l.Topics {
			topics[i] = web3.Hash(v)
		}
		parsed, err := tracer.eventEntry.ParseLog(&web3.Log{Address: web3.Address(l.Address), Topics: topics, Data: l.Data})
		if err != nil {
			return nil, fmt.Errorf("parse log failed err: %s, log info: %s", err, utils.JsonStr(l))
		}
		ret = append(ret, parsed)
	}
	res, err := json.Marshal(ret)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(res), nil
}

// Stop terminates execution of the tracer at the first opportune moment.
func (tracer *logTracer) Stop(err error) {
	tracer.reason = err
	atomic.StoreUint32(&tracer.interrupt, 1)
}

// make log instruction function
func makeLog(size int) func(pc uint64, scope *vm.ScopeContext) *types.Log {
	return func(pc uint64, scope *vm.ScopeContext) *types.Log {
		topics := make([]common.Hash, size)
		stack := scope.Stack
		mStart, mSize := stack.Back(0), stack.Back(1)
		for i := 0; i < size; i++ {
			addr := stack.Back(i + 2)
			topics[i] = addr.Bytes32()
		}

		d := scope.Memory.GetCopy(int64(mStart.Uint64()), int64(mSize.Uint64()))
		return &types.Log{
			Address: scope.Contract.Address(),
			Topics:  topics,
			Data:    d,
		}
	}
}
