package main

import "github.com/ethereum/go-ethereum/console"

func injectApiSignature(c *console.Console) {
	c.Evaluate(`
			if (l2) {  
				if(l2.getBatch){
					l2.getBatch.toString = function() { return "function (batchIndex: uint64, useDetail: bool)" }
				}
				if(l2.getBatchIndex){
					l2.getBatchIndex.toString = function() { return "function (blockNumber: uint64)" }
				}
				if(l2.getBatchState){
					l2.getBatchState.toString = function() { return "function (batchIndex: uint64)" }
				}
				if(l2.getEnqueuedTxs){
					l2.getEnqueuedTxs.toString = function() { return "function (queueStart: uint64, queueNum: uint64)" }
				}
				if(l2.getL2MMRProof){
					l2.getL2MMRProof.toString = function() { return "function (msgIndex: uint64, size: uint64)" }
				}
				if(l2.getL2RelayMsgParams){
					l2.getL2RelayMsgParams.toString = function() { return "function (msgIndex: uint64)" }
				}
				if(l2.getPendingTxBatches){
					l2.getPendingTxBatches.toString = function() { return "function ()" }
				}
				if(l2.globalInfo){
					l2.globalInfo.toString = function() { return "function ()" }
				}
				if(l2.getState){
					l2.getState.toString = function() { return "function (batchIndex: uint64)" }
				}
				if(l2.inputBatchNumber){
					l2.inputBatchNumber.toString = function() { return "function ()" }
				}
				if(l2.stateBatchNumber){
					l2.stateBatchNumber.toString = function() { return "function ()" }
				}
			}

			if (web3) {
				if (web3.sha3) {
					web3.sha3.toString = function() { return "function (data: bytes)" }
				}
				
			}

			if (eth) {
				if (eth.getBalance) {
					eth.getBalance.toString = function() { return "function (address: bytes20, blockNumber: uint or latest earliest pending)" }
				}
				if (eth.getStorageAt) {
					eth.getStorageAt.toString = function() { return "function (storageAddress: bytes20, storagePos: uint, blockNumber: uint or latest earliest pending)" }
				}
				if (eth.getTransactionCount) {
					eth.getTransactionCount.toString = function() { return "function (address: bytes20, blockNumber: uint or latest earliest pending)" }
				}
				if (eth.getBlockTransactionCount) {
					eth.getBlockTransactionCount.toString = function() { return "function (args: string[])" }
				}
				if (eth.getCode) {
					eth.getCode.toString = function() { return "function (address: bytes20, blockNumber: uint or latest earliest pending)" }
				}
				if (eth.sign) {
					eth.sign.toString = function() { return "function (address: bytes20, data: bytes)" }
				}
				if (eth.signTransaction) {
					eth.signTransaction.toString = function() { return "function (transaction: object(from: bytes20, to: bytes20, gas: uint, gasPrice: uint, value: uint, data: bytes, nonce: uint))" }
				}
				if (eth.sendTransaction) {
					eth.sendTransaction.toString = function() { return "function (transaction: object(from: bytes20, to: bytes20, gas: uint, gasPrice: uint, value: uint, data: bytes, nonce: uint))" }
				}
				if (eth.sendRawTransaction) {
					eth.sendRawTransaction.toString = function() { return "function (signedTransactionData: bytes)" }
				}
				if (eth.call) {
					eth.call.toString = function() { return "function ((transaction: object(from: bytes20, to: bytes20, gas: uint, gasPrice: uint, value: uint, data: bytes, nonce: uint), blockNumber: uint or latest earliest pending)" }
				}
				if (eth.estimateGas) {
					eth.estimateGas.toString = function() { return "function (transaction: object(from: bytes20, to: bytes20, gas: uint, gasPrice: uint, value: uint, data: bytes, nonce: uint), blockNumber: uint or latest earliest pending)" }
				}
				if (eth.getBlockByHash) {
					eth.getBlockByHash.toString = function() { return "function (blockHash: bytes32, isFullTransaction: bool)" }
				}
				if (eth.getBlockByNumber) {
					eth.getBlockByNumber.toString = function() { return "function (blockNumber: uint or latest earliest pending, isFullTransaction: bool)" }
				}
				if (eth.getTransactionReceipt) {
					eth.getTransactionReceipt.toString = function() { return "function (transactionHash: bytes32)" }
				}
				if (eth.compile && eth.compile.solidity) {
					eth.compile.solidity.toString = function() { return "function (sourceCode: string)" }
				}
				if (eth.compile && eth.compile.lll) {
					eth.compile.lll.toString = function() { return "function (sourceCode: string)" }
				}
				if (eth.compile && eth.compile.serpent) {
					eth.compile.serpent.toString = function() { return "function (sourceCode: string)" }
				}
				if (eth.getWork) {
					eth.getWork.toString = function() { return "function ()" }
				}
				if (eth.submitWork) {
					eth.submitWork.toString = function() { return "function (nonceFound: bytes8, headersPowHash: bytes32, mixDigest: bytes32)" }
				}
				if (eth.contract) {
					eth.contract.toString = function() { return "function (abi: bytes)" }
				}
				if (eth.createAccessList) {
					eth.createAccessList.toString = function() { return "function (transaction: object(from: bytes20, to: bytes20, gas: uint, gasPrice: uint, value: uint, data: bytes, nonce: uint))" }
				}
				if (eth.fillTransaction) {
					eth.fillTransaction.toString = function() { return "function (transaction: object(from: bytes20, to: bytes20, gas: uint, gasPrice: uint, value: uint, data: bytes, nonce: uint))" }
				}
				if (eth.getBlock) {
					eth.getBlock.toString = function() { return "function (blockNumberOrBlockHash: uint latest earliest pending or bytes)" }
				}
				if (eth.getBlockUncleCount) {
					eth.getBlockUncleCount.toString = function() { return "function (blockNumberOrBlockHash: uint latest earliest pending or bytes)" }
				}
				if (eth.getHeaderByHash) {
					eth.getHeaderByHash.toString = function() { return "function (BlockHash: bytes)" }
				}
				if (eth.getHeaderByNumber) {
					eth.getHeaderByNumber.toString = function() { return "function (BlockNumber: uint)" }
				}
				if (eth.getProof) {
					eth.getProof.toString = function() { return "function (address: bytes20, key: bytes[], blockNumberOrBlockHash: uint latest earliest pending or bytes)" }
				}
				if (eth.getRawTransaction) {
					eth.getRawTransaction.toString = function() { return "function (transactionHash: bytes)" }
				}
				if (eth.getTransaction) {
					eth.getTransaction.toString = function() { return "function (transactionHash: bytes)" }
				}
				if (eth.getTransactionFromBlock) {
					eth.getTransactionFromBlock.toString = function() { return "function (blockNumberOrBlockHash: uint latest earliest pending or bytes, txIndex: uint, callback)" }
				}
				if (eth.getUncle) {
					eth.getUncle.toString = function() { return "function (blockNumberOrBlockHash: uint latest earliest pending or bytes, uncleIndex: uint, isReturnTransactionObjects: bool, callback)" }
				}
				if (eth.getWork) {
					eth.getWork.toString = function() { return "function (callback)" }
				}
				if (eth.iban) {
					eth.iban.toString = function() { return "function (ethAddress: bytes20)" }
				}
				if (eth.resend) {
					eth.resend.toString = function() { return "function (transaction: object(from: bytes20, to: bytes20, gas: uint, gasPrice: uint, value: uint, data: bytes, nonce: uint), gasPrice: uint, gasLimit: uint)" }
				}
			}


			if (debug) {  
				if (debug.accountRange) {
					debug.accountRange.toString = function() { return "function (blockNumberOrBlockHash: uint latest earliest pending or bytes, start: string, maxResults: uint, skip: uint, reverse: bool)" }
				}
				if (debug.backtraceAt) {
					debug.backtraceAt.toString = function() { return "function ((filename:line): string)" }
				}
				if (debug.blockProfile) {
					debug.blockProfile.toString = function() { return "function (filename: string, time: uint)" }
				}
				if (debug.chaindbCompact) {
					debug.chaindbCompact.toString = function() { return "function ()" }
				}
				if (debug.chaindbProperty) {
					debug.chaindbProperty.toString = function() { return "function (leveldb properties: string)" }
				}
				if (debug.cpuProfile) {
					debug.cpuProfile.toString = function() { return "function (filename: string, time: uint)" }
				}
				if (debug.dumpBlock) {
					debug.dumpBlock.toString = function() { return "function (blockNumber: uint or latest earliest pending)" }
				}
				if (debug.freeOSMemory) {
					debug.freeOSMemory.toString = function() { return "function ()" }
				}
				if (debug.freezeClient) {
					debug.freezeClient.toString = function() { return "function ()" }
				}
				if (debug.gcStats) {
					debug.gcStats.toString = function() { return "function ()" }
				}
				if (debug.getAccessibleState) {
					debug.getAccessibleState.toString = function() { return "function (from: uint, to: uint)" }
				}
				if (debug.getBadBlocks) {
					debug.getBadBlocks.toString = function() { return "function ()" }
				}
				if (debug.getBlockRlp) {
					debug.getBlockRlp.toString = function() { return "function (blockNumberOrBlockHash: uint latest earliest pending or bytes)" }
				}
				if (debug.getHeaderRlp) {
					debug.getHeaderRlp.toString = function() { return "function (blockNumber: uint)" }
				}
				if (debug.getModifiedAccountsByHash) {
					debug.getModifiedAccountsByHash.toString = function() { return "function (startHash: bytes, endHash: bytes)" }
				}
				if (debug.getModifiedAccountsByNumber) {
					debug.getModifiedAccountsByNumber.toString = function() { return "function (startNum: uint64, endNum: uint64)" }
				}
				if (debug.goTrace) {
					debug.goTrace.toString = function() { return "function (filename: string, runningTime: uint)" }
				}
				if (debug.intermediateRoots) {
					debug.intermediateRoots.toString = function() { return "function (blockHash: string)" }
				}
				if (debug.mutexProfile) {
					debug.mutexProfile.toString = function() { return "function (filename: string, nsec: uint)" }
				}
				if (debug.preimage) {
					debug.preimage.toString = function() { return "function (hash: string)" }
				}
				if (debug.printBlock) {
					debug.printBlock.toString = function() { return "function (blockNumber: uint)" }
				}
				if (debug.seedHash) {
					debug.seedHash.toString = function() { return "function (blockNumber: uint)" }
				}
				if (debug.setBlockProfileRate) {
					debug.setBlockProfileRate.toString = function() { return "function (rate: int)" }
				}
				if (debug.setGCPercent) {
					debug.setGCPercent.toString = function() { return "function (v: int)" }
				}
				if (debug.setHead) {
					debug.setHead.toString = function() { return "function (blockNumber: string)" }
				}
				if (debug.setMutexProfileFraction) {
					debug.setMutexProfileFraction.toString = function() { return "function (rate: int)" }
				}
				if (debug.stacks) {
					debug.stacks.toString = function() { return "function ()" }
				}
				if (debug.standardTraceBadBlockToFile) {
					debug.standardTraceBadBlockToFile.toString = function() { return "function (blockHash: string, optionObject(disableStorage: bool, disableMemory: bool, disableStack: bool, fullStorage: bool))" }
				}
				if (debug.standardTraceBlockToFile) {
					debug.standardTraceBlockToFile.toString = function() { return "function (blockHash: string, optionObject(disableStorage: bool, disableMemory: bool, disableStack: bool, fullStorage: bool))" }
				}
				if (debug.startCPUProfile) {
					debug.startCPUProfile.toString = function() { return "function (filename: string)" }
				}
				if (debug.startGoTrace) {
					debug.startGoTrace.toString = function() { return "function (filename: string)" }
				}
				if (debug.stopCPUProfile) {
					debug.stopCPUProfile.toString = function() { return "function ()" }
				}
				if (debug.stopGoTrace) {
					debug.stopGoTrace.toString = function() { return "function ()" }
				}
				if (debug.storageRangeAt) {
					debug.storageRangeAt.toString = function() { return "function (blockHash: string, transactionIdx: uint, contractAddress: bytes20, keyStart: string, maxResult: uint)" }
				}
				if (debug.testSignCliqueBlock) {
					debug.testSignCliqueBlock.toString = function() { return "function (sigAddress: bytes20, blockNumber: uint)" }
				}
				if (debug.traceBadBlock) {
					debug.traceBadBlock.toString = function() { return "function (blockHash: string, optionObject: (disableStorage: bool, disableStack: bool, enableMemory: bool, enableReturnData: bool, tracer: string, timeOut: string, tracerConfig: object(onlyTopCall: bool, withLog: bool)))" }
				}
				if (debug.traceBlock) {
					debug.traceBlock.toString = function() { return "function (blockRlp: string, optionObject: (disableStorage: bool, disableStack: bool, enableMemory: bool, enableReturnData: bool, tracer: string, timeOut: string, tracerConfig: object(onlyTopCall: bool, withLog: bool)))" }
				}
				if (debug.traceBlockByHash) {
					debug.traceBlockByHash.toString = function() { return "function (blockHash: string, optionObject: (disableStorage: bool, disableStack: bool, enableMemory: bool, enableReturnData: bool, tracer: string, timeOut: string, tracerConfig: object(onlyTopCall: bool, withLog: bool)))" }
				}
				if (debug.traceBlockByNumber) {
					debug.traceBlockByNumber.toString = function() { return "function (blockNumber: uint, optionObject: (disableStorage: bool, disableStack: bool, enableMemory: bool, enableReturnData: bool, tracer: string, timeOut: string, tracerConfig: object(onlyTopCall: bool, withLog: bool)))" }
				}
				if (debug.traceBlockFromFile) {
					debug.traceBlockFromFile.toString = function() { return "function (fileName: string, optionObject: (disableStorage: bool, disableStack: bool, enableMemory: bool, enableReturnData: bool, tracer: string, timeOut: string, tracerConfig: object(onlyTopCall: bool, withLog: bool)))" }
				}
				if (debug.traceCall) {
					debug.traceCall.toString = function() { return "function ((transaction: object(from: bytes20, to: bytes20, gas: uint, gasPrice: uint, value: uint, data: bytes, nonce: uint), blockNrOrHash: uint latest earliest pending or bytes, optionObject: (disableStorage: bool, disableStack: bool, enableMemory: bool, enableReturnData: bool, tracer: string, timeOut: string, tracerConfig: object(onlyTopCall: bool, withLog: bool)))" }
				}
				if (debug.traceTransaction) {
					debug.traceTransaction.toString = function() { return "function (transactionHash: bytes, optionObject: (disableStorage: bool, disableStack: bool, enableMemory: bool, enableReturnData: bool, tracer: string, timeOut: string, tracerConfig: object(onlyTopCall: bool, withLog: bool)))" }
				}
				if (debug.verbosity) {
					debug.verbosity.toString = function() { return "function (logLevel: int)" }
				}
				if (debug.vmodule) {
					debug.vmodule.toString = function() { return "function (logPattern: string)" }
				}
				if (debug.writeBlockProfile) {
					debug.writeBlockProfile.toString = function() { return "function (file: string, seconds: uint)" }
				}
				if (debug.writeMemProfile) {
					debug.writeMemProfile.toString = function() { return "function (file: string)" }
				}
				if (debug.writeMutexProfile) {
					debug.writeMutexProfile.toString = function() { return "function (file: string)" }
				}
				if (debug.getReadStorageProofAtBlock) {
					debug.getReadStorageProofAtBlock.toString = function() { return "function (blockNumber: uint64)" }
				}
			}
	`)
}
