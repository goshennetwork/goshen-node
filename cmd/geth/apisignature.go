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
		debug.getReadStorageProofAtBlock.toString = function() { return "function (blockNumber: uint64)" }
	`)
}
