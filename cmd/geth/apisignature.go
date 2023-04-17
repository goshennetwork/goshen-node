package main

import "github.com/ethereum/go-ethereum/console"

func injectApiSignature(c *console.Console) {
	c.Evaluate(`
			l2.getBatch.toString = function() { return "function (batchIndex: uint64, useDetail: bool)" }
			l2.getBatchIndex.toString = function() { return "function (blockNumber: uint64)" }
			l2.getBatchState.toString = function() { return "function (batchIndex: uint64)" }
			l2.getEnqueuedTxs.toString = function() { return "function (queueStart: uint64, queueNum: uint64)" }
			l2.getL2MMRProof.toString = function() { return "function (msgIndex: uint64, size: uint64)" }
			l2.getL2RelayMsgParams.toString = function() { return "function (msgIndex: uint64)" }
			l2.getPendingTxBatches.toString = function() { return "function ()" }
			l2.globalInfo.toString = function() { return "function ()" }
			l2.getState.toString = function() { return "function (batchIndex: uint64)" }
			l2.inputBatchNumber.toString = function() { return "function ()" }
			l2.stateBatchNumber.toString = function() { return "function ()" }

			web3.sha3.toString = function() { return "function (data: bytes)" }
			eth.getBalance.toString = function() { return "function (address: bytes20, blockNumber: uint or latest earliest pending)" }
			eth.getStorageAt.toString = function() { return "function (storageAddress: bytes20, storagePos: uint, blockNumber: uint or latest earliest pending)" }
			eth.getTransactionCount.toString = function() { return "function (address: bytes20, blockNumber: uint or latest earliest pending)" }
			eth.getBlockTransactionCount.toString = function() { return "function (args: string[])" }
			eth.getCode.toString = function() { return "function (address: bytes20, blockNumber: uint or latest earliest pending)" }
			eth.sign.toString = function() { return "function (address: bytes20, data: bytes)" }
			eth.signTransaction.toString = function() { return "function (transaction: object(from: bytes20, to: bytes20, gas: uint, gasPrice: uint, value: uint, data: bytes, nonce: uint))" }
			eth.sendTransaction.toString = function() { return "function (transaction: object(from: bytes20, to: bytes20, gas: uint, gasPrice: uint, value: uint, data: bytes, nonce: uint))" }
			eth.sendRawTransaction.toString = function() { return "function (signedTransactionData: bytes)" }
			eth.call.toString = function() { return "function ((transaction: object(from: bytes20, to: bytes20, gas: uint, gasPrice: uint, value: uint, data: bytes, nonce: uint), blockNumber: uint or latest earliest pending)" }
			eth.estimateGas.toString = function() { return "function (transaction: object(from: bytes20, to: bytes20, gas: uint, gasPrice: uint, value: uint, data: bytes, nonce: uint), blockNumber: uint or latest earliest pending)" }
			eth.getBlockByHash.toString = function() { return "function (blockHash: bytes32, isFullTransaction: bool)" }
			eth.getBlockByNumber.toString = function() { return "function (blockNumber: uint or latest earliest pending, isFullTransaction: bool)" }
			eth.getTransactionReceipt.toString = function() { return "function (transactionHash: bytes32)" }
			eth.compile.solidity.toString = function() { return "function (sourceCode: string)" }
			eth.compile.lll.toString = function() { return "function (sourceCode: string)" }
			eth.compile.serpent.toString = function() { return "function (sourceCode: string)" }
			eth.getWork.toString = function() { return "function ()" }
			eth.submitWork.toString = function() { return "function (nonceFound: bytes8, headersPowHash: bytes32, mixDigest: bytes32)" }
			eth.contract.toString = function() { return "function (abi: bytes)" }
			eth.createAccessList.toString = function() { return "function (transaction: object(from: bytes20, to: bytes20, gas: uint, gasPrice: uint, value: uint, data: bytes, nonce: uint))" }
			eth.fillTransaction.toString = function() { return "function (transaction: object(from: bytes20, to: bytes20, gas: uint, gasPrice: uint, value: uint, data: bytes, nonce: uint))" }
			eth.getBlock.toString = function() { return "function (blockNumberOrBlockHash: uint latest earliest pending or bytes)" }
			eth.getBlockUncleCount.toString = function() { return "function (blockNumberOrBlockHash: uint latest earliest pending or bytes)" }
			eth.getHeaderByHash.toString = function() { return "function (BlockHash: bytes)" }
			eth.getHeaderByNumber.toString = function() { return "function (BlockNumber: uint)" }
			eth.getProof.toString = function() { return "function (address: bytes20, key: bytes[], blockNumberOrBlockHash: uint latest earliest pending or bytes)" }
			eth.getRawTransaction.toString = function() { return "function (transactionHash: bytes)" }
			eth.getTransaction.toString = function() { return "function (transactionHash: bytes)" }
			eth.getTransactionFromBlock.toString = function() { return "function (blockNumberOrBlockHash: uint latest earliest pending or bytes, txIndex: uint, callback)" }
			eth.getUncle.toString = function() { return "function (blockNumberOrBlockHash: uint latest earliest pending or bytes, uncleIndex: uint, isReturnTransactionObjects: bool, callback)" }
			eth.getWork.toString = function() { return "function (callback)" }
			eth.iban.toString = function() { return "function (ethAddress: bytes20)" }
			eth.resend.toString = function() { return "function (transaction: object(from: bytes20, to: bytes20, gas: uint, gasPrice: uint, value: uint, data: bytes, nonce: uint), gasPrice: uint, gasLimit: uint)" }
			
			
			debug.accountRange.toString = function() { return "function (blockNumberOrBlockHash: uint latest earliest pending or bytes, start: string, maxResults: uint, skip: uint, reverse: bool)" }
			debug.backtraceAt.toString = function() { return "function ((filename:line): string)" }
			debug.blockProfile.toString = function() { return "function (filename: string, time: uint)" }
			debug.chaindbCompact.toString = function() { return "function ()" }
			debug.chaindbProperty.toString = function() { return "function (leveldb properties: string)" }
			debug.cpuProfile.toString = function() { return "function (filename: string, time: uint)" }
			debug.dumpBlock.toString = function() { return "function (blockNumber: uint or latest earliest pending)" }
			debug.freeOSMemory.toString = function() { return "function ()" }
			debug.freezeClient.toString = function() { return "function ()" }
			debug.gcStats.toString = function() { return "function ()" }
			debug.getAccessibleState.toString = function() { return "function (from: uint, to: uint)" }
			debug.getBadBlocks.toString = function() { return "function ()" }
			debug.getBlockRlp.toString = function() { return "function (blockNumberOrBlockHash: uint latest earliest pending or bytes)" }
			debug.getHeaderRlp.toString = function() { return "function (blockNumber: uint)" }
			debug.getModifiedAccountsByHash.toString = function() { return "function (startHash: bytes, endHash: bytes)" }
			debug.getModifiedAccountsByNumber.toString = function() { return "function (startNum: uint64, endNum: uint64)" }
			debug.goTrace.toString = function() { return "function (filename: string, runningTime: uint)" }
			debug.intermediateRoots.toString = function() { return "function (blockHash: string)" }
			debug.mutexProfile.toString = function() { return "function (filename: string, nsec: uint)" }
			debug.preimage.toString = function() { return "function (hash: string)" }
			debug.printBlock.toString = function() { return "function (blockNumber: uint)" }
			debug.seedHash.toString = function() { return "function (blockNumber: uint)" }
			debug.setBlockProfileRate.toString = function() { return "function (rate: int)" }
			debug.setGCPercent.toString = function() { return "function (v: int)" }
			debug.setHead.toString = function() { return "function (blockNumber: string)" }
			debug.setMutexProfileFraction.toString = function() { return "function (rate: int)" }
			debug.stacks.toString = function() { return "function ()" }
			debug.standardTraceBadBlockToFile.toString = function() { return "function (blockHash: string, optionObject(disableStorage: bool, disableMemory: bool, disableStack: bool, fullStorage: bool))" }
			debug.standardTraceBlockToFile.toString = function() { return "function (blockHash: string, optionObject(disableStorage: bool, disableMemory: bool, disableStack: bool, fullStorage: bool))" }
			debug.startCPUProfile.toString = function() { return "function (filename: string)" }
			debug.startGoTrace.toString = function() { return "function (filename: string)" }
			debug.stopCPUProfile.toString = function() { return "function ()" }
			debug.stopGoTrace.toString = function() { return "function ()" }
			debug.storageRangeAt.toString = function() { return "function (blockHash: string, transactionIdx: uint, contractAddress: bytes20, keyStart: string, maxResult: uint)" }
			debug.testSignCliqueBlock.toString = function() { return "function (sigAddress: bytes20, blockNumber: uint)" }
			debug.traceBadBlock.toString = function() { return "function (blockHash: string, optionObject: (disableStorage: bool, disableStack: bool, enableMemory: bool, enableReturnData: bool, tracer: string, timeOut: string, tracerConfig: object(onlyTopCall: bool, withLog: bool)))" }
			debug.traceBlock.toString = function() { return "function (blockRlp: string, optionObject: (disableStorage: bool, disableStack: bool, enableMemory: bool, enableReturnData: bool, tracer: string, timeOut: string, tracerConfig: object(onlyTopCall: bool, withLog: bool)))" }
			debug.traceBlockByHash.toString = function() { return "function (blockHash: string, optionObject: (disableStorage: bool, disableStack: bool, enableMemory: bool, enableReturnData: bool, tracer: string, timeOut: string, tracerConfig: object(onlyTopCall: bool, withLog: bool)))" }
			debug.traceBlockByNumber.toString = function() { return "function (blockNumber: uint, optionObject: (disableStorage: bool, disableStack: bool, enableMemory: bool, enableReturnData: bool, tracer: string, timeOut: string, tracerConfig: object(onlyTopCall: bool, withLog: bool)))" }
			debug.traceBlockFromFile.toString = function() { return "function (fileName: string, optionObject: (disableStorage: bool, disableStack: bool, enableMemory: bool, enableReturnData: bool, tracer: string, timeOut: string, tracerConfig: object(onlyTopCall: bool, withLog: bool)))" }
			debug.traceCall.toString = function() { return "function ((transaction: object(from: bytes20, to: bytes20, gas: uint, gasPrice: uint, value: uint, data: bytes, nonce: uint), blockNrOrHash: uint latest earliest pending or bytes, optionObject: (disableStorage: bool, disableStack: bool, enableMemory: bool, enableReturnData: bool, tracer: string, timeOut: string, tracerConfig: object(onlyTopCall: bool, withLog: bool)))" }
			debug.traceTransaction.toString = function() { return "function (transactionHash: bytes, optionObject: (disableStorage: bool, disableStack: bool, enableMemory: bool, enableReturnData: bool, tracer: string, timeOut: string, tracerConfig: object(onlyTopCall: bool, withLog: bool)))" }
			debug.verbosity.toString = function() { return "function (logLevel: int)" }
			debug.vmodule.toString = function() { return "function (logPattern: string)" }
			debug.writeBlockProfile.toString = function() { return "function (file: string, seconds: uint)" }
			debug.writeMemProfile.toString = function() { return "function (file: string)" }
			debug.writeMutexProfile.toString = function() { return "function (file: string)" }
	`)
}
