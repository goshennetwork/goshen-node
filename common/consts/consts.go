package consts

import (
	"flag"
	"math"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

const InitialEnqueueTxNonce = 1 << 63
const MaxSenderNonce = 1 << 62
const MaxL1TxSize = 128 * 1024 // defined in tx_pool.go of go-ethereum
const MaxL2TxSize = 32 * 1024  // limit L2 tx to 32KB to ensure it can be submitted to L1

// Compute maximal data size for transactions (lower bound).
//
// It is assumed the fields in the transaction (except of the data) are:
//   - nonce     <= 32 bytes
//   - gasPrice  <= 32 bytes
//   - gasLimit  <= 32 bytes
//   - recipient == 20 bytes
//   - value     <= 32 bytes
//   - signature == 65 bytes
// All those fields are summed up to at most 213 bytes.
const TxBaseSize = 213
const MaxL1TxDataSize = MaxL1TxSize - TxBaseSize
const MaxRollupInputBatchSize = MaxL1TxSize*48/128 - TxBaseSize // 48KB

const GenesisDifficulty = 1

var L1CrossLayerWitnessSender = common.HexToAddress("0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf")

var FeeCollector = common.HexToAddress("0xfee0000000000000000000000000000000000fee")

var IntrinsicGasFactor = func() uint64 {
	if IsTestintEnv() {
		return 1
	}

	return 100 // we want L2 gas price 1% of L1, except intrinsic gas
}()

// limit the execution gas: tx.Gas - tx.IntrinsicGas
var MaxTxExecGas = func() uint64 {
	if IsTestintEnv() {
		return math.MaxUint64
	}

	return 20000000
}()

func IsTestintEnv() bool {
	return flag.Lookup("test.v") != nil || strings.Contains(os.Args[0], "test")
}
