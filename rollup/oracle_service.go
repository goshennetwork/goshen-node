package rollup

import (
	"errors"
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common/consts"
	"github.com/laizy/web3"
)

type PriceOracleService struct {
	*RollupBackend
	l1price uint64
	quit    chan struct{}
	running uint32
}

var DefaultL1Price = web3.Gwei(100).Uint64()

func NewPriceOracleService(backend *RollupBackend) *PriceOracleService {
	return &PriceOracleService{backend, DefaultL1Price, make(chan struct{}, 1), 0}
}

func (self *PriceOracleService) Start() error {
	if !atomic.CompareAndSwapUint32(&self.running, 0, 1) {
		return errors.New("already running")
	}
	go self.run()
	return nil
}

func (self *PriceOracleService) Stop() error {
	if !atomic.CompareAndSwapUint32(&self.running, 1, 0) {
		return errors.New("already closed")
	}
	self.quit <- struct{}{}
	return nil
}

func (self *PriceOracleService) SetL1Price(price uint64) { atomic.StoreUint64(&self.l1price, price) }
func (self *PriceOracleService) L1Price() uint64         { return atomic.LoadUint64(&self.l1price) }
func (self *PriceOracleService) IsRunning() bool         { return atomic.LoadUint32(&self.running) == 1 }
func (self *PriceOracleService) L2Price(minPrice *big.Int) (*big.Int, error) {
	if !self.IsRunning() {
		return nil, errors.New("l1 gasPrice oracle not running")
	}

	l1price := self.L1Price()
	l2price := new(big.Int).Div(new(big.Int).SetUint64(l1price), new(big.Int).SetUint64(uint64(consts.IntrinsicGasFactor)))
	if l2price.Cmp(minPrice) < 0 {
		l2price.Set(minPrice) //0.1gwei
	}
	return l2price, nil
}
