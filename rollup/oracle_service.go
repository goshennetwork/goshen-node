package rollup

import (
	"errors"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common/consts"
	"github.com/ethereum/go-ethereum/log"
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

func (self *PriceOracleService) run() {
	var secondTime time.Duration = 12
	var minuteTime time.Duration = 10
	ticker := time.NewTicker(secondTime * time.Second)
	defer ticker.Stop()
	var l1gasPricesQue []uint64

	for {
		select {
		case <-self.quit:
			log.Warn("price oracle service", "info", "quiting")
			return
		case <-ticker.C:
			l1price, err := self.L1Client.Eth().GasPrice()
			if err != nil {
				log.Warn("get l1 gasprice", "err", err)
			} else {
				if len(l1gasPricesQue) == int(minuteTime*time.Minute/secondTime*time.Second) {
					l1gasPricesQue = l1gasPricesQue[1:]
					l1gasPricesQue = append(l1gasPricesQue, l1price)
				} else {
					l1gasPricesQue = append(l1gasPricesQue, l1price)
				}
				l1maxGasPrice := l1gasPricesQue[0]
				for _, price := range l1gasPricesQue {
					if price > l1maxGasPrice {
						l1maxGasPrice = price
					}
				}
				self.SetL1Price(l1maxGasPrice)
			}
		}
	}
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
