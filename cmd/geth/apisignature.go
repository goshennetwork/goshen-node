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


		if (eth) {
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
			eth.getBlockByHash.toString = function() { return "function (blockHash: string, isFullTransaction: bool)" }
			eth.getBlockByNumber.toString = function() { return "function (blockNumber: uint or latest earliest pending, isFullTransaction: bool)" }
			eth.getTransactionReceipt.toString = function() { return "function (transactionHash: string)" }
			eth.compile.solidity.toString = function() { return "function (sourceCode: string)" }
			eth.compile.lll.toString = function() { return "function (sourceCode: string)" }
			eth.compile.serpent.toString = function() { return "function (sourceCode: string)" }
			eth.getWork.toString = function() { return "function ()" }
			eth.submitWork.toString = function() { return "function (nonceFound: bytes8, headersPowHash: string, mixDigest: bytes32)" }
			eth.contract.toString = function() { return "function (abi: bytes)" }
			eth.createAccessList.toString = function() { return "function (transaction: object(from: bytes20, to: bytes20, gas: uint, gasPrice: uint, value: uint, data: bytes, nonce: uint))" }
			eth.fillTransaction.toString = function() { return "function (transaction: object(from: bytes20, to: bytes20, gas: uint, gasPrice: uint, value: uint, data: bytes, nonce: uint))" }
			eth.getBlock.toString = function() { return "function (blockNumberOrBlockHash: uint latest earliest pending or string)" }
			eth.getBlockUncleCount.toString = function() { return "function (blockNumberOrBlockHash: uint latest earliest pending or string)" }
			eth.getHeaderByHash.toString = function() { return "function (BlockHash: string)" }
			eth.getHeaderByNumber.toString = function() { return "function (BlockNumber: uint)" }
			eth.getProof.toString = function() { return "function (address: bytes20, key: bytes[], blockNumberOrBlockHash: uint latest earliest pending or string)" }
			eth.getRawTransaction.toString = function() { return "function (transactionHash: string)" }
			eth.getTransaction.toString = function() { return "function (transactionHash: string)" }
			eth.getTransactionFromBlock.toString = function() { return "function (blockNumberOrBlockHash: uint latest earliest pending or string, transactionIndex: uint, callback)" }
			eth.getUncle.toString = function() { return "function (blockNumberOrBlockHash: uint latest earliest pending or string, uncleIndex: uint, isReturnTransactionObjects: bool, callback)" }
			eth.getWork.toString = function() { return "function (callback)" }
			eth.iban.toString = function() { return "function (ethAddress: bytes20)" }
			eth.resend.toString = function() { return "function (transaction: object(from: bytes20, to: bytes20, gas: uint, gasPrice: uint, value: uint, data: bytes, nonce: uint), gasPrice: uint, gasLimit: uint)" }
		}


		if (debug) {  
			if (debug.accountRange) {
				debug.accountRange.toString = function() { return "function (blockNumberOrBlockHash: uint latest earliest pending or string, start: string, maxResults: uint, skip: uint, reverse: bool)" }
			}
			if (debug.backtraceAt) {
				debug.backtraceAt.toString = function() { return "function ((file:line): string)" }
			}
			if (debug.blockProfile) {
				debug.blockProfile.toString = function() { return "function (file: string, time: uint)" }
			}
			if (debug.chaindbCompact) {
				debug.chaindbCompact.toString = function() { return "function ()" }
			}
			if (debug.chaindbProperty) {
				debug.chaindbProperty.toString = function() { return "function (leveldb properties: string)" }
			}
			if (debug.cpuProfile) {
				debug.cpuProfile.toString = function() { return "function (file: string, time: uint)" }
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
				debug.getBlockRlp.toString = function() { return "function (blockNumberOrBlockHash: uint latest earliest pending or string)" }
			}
			if (debug.getHeaderRlp) {
				debug.getHeaderRlp.toString = function() { return "function (blockNumber: uint)" }
			}
			if (debug.getModifiedAccountsByHash) {
				debug.getModifiedAccountsByHash.toString = function() { return "function (startHash: string, endHash: string)" }
			}
			if (debug.getModifiedAccountsByNumber) {
				debug.getModifiedAccountsByNumber.toString = function() { return "function (startNum: uint64, endNum: uint64)" }
			}
			if (debug.goTrace) {
				debug.goTrace.toString = function() { return "function (file: string, runningTime: uint)" }
			}
			if (debug.intermediateRoots) {
				debug.intermediateRoots.toString = function() { return "function (blockHash: string)" }
			}
			if (debug.mutexProfile) {
				debug.mutexProfile.toString = function() { return "function (file: string, nsec: uint)" }
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
				debug.standardTraceBadBlockToFile.toString = function() { return "function (blockHash: string, optionObject: object(disableStorage: bool, disableMemory: bool, disableStack: bool, fullStorage: bool))" }
			}
			if (debug.standardTraceBlockToFile) {
				debug.standardTraceBlockToFile.toString = function() { return "function (blockHash: string, optionObject: object(disableStorage: bool, disableMemory: bool, disableStack: bool, fullStorage: bool))" }
			}
			if (debug.startCPUProfile) {
				debug.startCPUProfile.toString = function() { return "function (file: string)" }
			}
			if (debug.startGoTrace) {
				debug.startGoTrace.toString = function() { return "function (file: string)" }
			}
			if (debug.stopCPUProfile) {
				debug.stopCPUProfile.toString = function() { return "function ()" }
			}
			if (debug.stopGoTrace) {
				debug.stopGoTrace.toString = function() { return "function ()" }
			}
			if (debug.storageRangeAt) {
				debug.storageRangeAt.toString = function() { return "function (blockHash: string, transactionIndex: uint, contractAddress: bytes20, keyStart: string, maxResult: uint)" }
			}
			if (debug.testSignCliqueBlock) {
				debug.testSignCliqueBlock.toString = function() { return "function (sigAddress: bytes20, blockNumber: uint)" }
			}
			if (debug.traceBadBlock) {
				debug.traceBadBlock.toString = function() { return "function (blockHash: string, optionObject: object(disableStorage: bool, disableStack: bool, enableMemory: bool, enableReturnData: bool, tracer: string, timeOut: string, tracerConfig: object(onlyTopCall: bool, withLog: bool)))" }
			}
			if (debug.traceBlock) {
				debug.traceBlock.toString = function() { return "function (blockRlp: string, optionObject: object(disableStorage: bool, disableStack: bool, enableMemory: bool, enableReturnData: bool, tracer: string, timeOut: string, tracerConfig: object(onlyTopCall: bool, withLog: bool)))" }
			}
			if (debug.traceBlockByHash) {
				debug.traceBlockByHash.toString = function() { return "function (blockHash: string, optionObject: object(disableStorage: bool, disableStack: bool, enableMemory: bool, enableReturnData: bool, tracer: string, timeOut: string, tracerConfig: object(onlyTopCall: bool, withLog: bool)))" }
			}
			if (debug.traceBlockByNumber) {
				debug.traceBlockByNumber.toString = function() { return "function (blockNumber: uint, optionObject: object(disableStorage: bool, disableStack: bool, enableMemory: bool, enableReturnData: bool, tracer: string, timeOut: string, tracerConfig: object(onlyTopCall: bool, withLog: bool)))" }
			}
			if (debug.traceBlockFromFile) {
				debug.traceBlockFromFile.toString = function() { return "function (file: string, optionObject: object(disableStorage: bool, disableStack: bool, enableMemory: bool, enableReturnData: bool, tracer: string, timeOut: string, tracerConfig: object(onlyTopCall: bool, withLog: bool)))" }
			}
			if (debug.traceCall) {
				debug.traceCall.toString = function() { return "function ((transaction: object(from: bytes20, to: bytes20, gas: uint, gasPrice: uint, value: uint, data: bytes, nonce: uint), blockNumberOrBlockHash: uint latest earliest pending or string, optionObject: (disableStorage: bool, disableStack: bool, enableMemory: bool, enableReturnData: bool, tracer: string, timeOut: string, tracerConfig: object(onlyTopCall: bool, withLog: bool)))" }
			}
			if (debug.traceTransaction) {
				debug.traceTransaction.toString = function() { return "function (transactionHash: string, optionObject: object(disableStorage: bool, disableStack: bool, enableMemory: bool, enableReturnData: bool, tracer: string, timeOut: string, tracerConfig: object(onlyTopCall: bool, withLog: bool)))" }
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


		if (miner) {
			if (miner.getHashrate) {
				miner.getHashrate.toString = function() { return "function ()" }
			}
			if (miner.setEtherbase) {
				miner.setEtherbase.toString = function() { return "function (address: bytes20)" }
			}
			if (miner.setExtra) {
				miner.setExtra.toString = function() { return "function (extraData: bytes32)" }
			}
			if (miner.setGasLimit) {
				miner.setGasLimit.toString = function() { return "function (gasLimit: uint)" }
			}
			if (miner.setGasPrice) {
				miner.setGasPrice.toString = function() { return "function (gasPrice: uint)" }
			}
			if (miner.setRecommitInterval) {
				miner.setRecommitInterval.toString = function() { return "function (timeInterval: uint)" }
			}
			if (miner.start) {
				miner.start.toString = function() { return "function (threadsNum: unit)" }
			}
			if (miner.stop) {
				miner.stop.toString = function() { return "function ()" }
			}
		}


		if (personal) {
			if (personal.deriveAccount) {
				personal.deriveAccount.toString = function() { return "function (seed: string, derivationPath: string)" }
			}
			if (personal.ecRecover) {
				personal.ecRecover.toString = function() { return "function (data: string, signature: string)" }
			}
			if (personal.getListAccounts) {
				personal.getListAccounts.toString = function() { return "function (callback)" }
			}
			if (personal.getListWallets) {
				personal.getListWallets.toString = function() { return "function (callback)" }
			}
			if (personal.importRawKey) {
				personal.importRawKey.toString = function() { return "function (privateKey: string, passphrase: string)" }
			}
			if (personal.initializeWallet) {
				personal.initializeWallet.toString = function() { return "function (url: string)" }
			}
			if (personal.lockAccount) {
				personal.lockAccount.toString = function() { return "function (address: bytes20)" }
			}
			if (personal.newAccount) {
				personal.newAccount.toString = function() { return "function (password: string)" }
			}
			if (personal.openWallet) {
				personal.openWallet.toString = function() { return "function (url: string, passphrase: string)" }
			}
			if (personal.sendTransaction) {
				personal.sendTransaction.toString = function() { return "function (transaction: object(from: bytes20, to: bytes20, gas: uint, gasPrice: uint, value: uint, data: bytes, nonce: uint), passphrase: string)" }
			}
			if (personal.sign) {
				personal.sign.toString = function() { return "function (dataToSign: string, accountAddress: bytes20, passphrase: string)" }
			}
			if (personal.signTransaction) {
				personal.signTransaction.toString = function() { return "function (transaction: object(from: bytes20, to: bytes20, gas: uint, gasPrice: uint, value: uint, data: bytes, nonce: uint), passphrase: string)" }
			}
			if (personal.unlockAccount) {
				personal.unlockAccount.toString = function() { return "function (accountAddress: bytes20, passphrase: string, time: uint)" }
			}
			if (personal.unpair) {
				personal.unpair.toString = function() { return "function (url: string, pin: string)" }
			}
		}


		if (utils) {
			if (utils.saveString) {
				utils.saveString.toString = function() { return "function (data: string, fileName: string)" }
			}
		}
		

		if (web3) {
			if (web3.shh) {
				if (web3.shh.addPrivateKey) {
					web3.shh.addPrivateKey.toString = function() { return "function (privateKey: string, callback)" }
				}
				if (web3.shh.addSymKey) {
					web3.shh.addSymKey.toString = function() { return "function (symKey: string, callback)" }
				}
				if (web3.shh.deleteKeyPair) {
					web3.shh.deleteKeyPair.toString = function() { return "function (id: string, callback)" }
				}
				if (web3.shh.deleteSymKey) {
					web3.shh.deleteSymKey.toString = function() { return "function (symKeyID: string, callback)" }
				}
				if (web3.shh.generateSymKeyFromPassword) {
					web3.shh.generateSymKeyFromPassword.toString = function() { return "function (password: string, callback)" }
				}
				if (web3.shh.getPrivateKey) {
					web3.shh.getPrivateKey.toString = function() { return "function (id: string, callback)" }
				}
				if (web3.shh.getPublicKey) {
					web3.shh.getPublicKey.toString = function() { return "function (id: string, callback)" }
				}
				if (web3.shh.getSymKey) {
					web3.shh.getSymKey.toString = function() { return "function (id: string, callback)" }
				}
				if (web3.shh.hasKeyPair) {
					web3.shh.hasKeyPair.toString = function() { return "function (id: string, callback)" }
				}
				if (web3.shh.hasSymKey) {
					web3.shh.hasSymKey.toString = function() { return "function (id: string, callback)" }
				}
				if (web3.shh.markTrustedPeer) {
					web3.shh.markTrustedPeer.toString = function() { return "function (enode: string, callback)" }
				}
				if (web3.shh.newKeyPair) {
					web3.shh.newKeyPair.toString = function() { return "function (callback)" }
				}
				if (web3.shh.newMessageFilter) {
					web3.shh.newMessageFilter.toString = function() { return "function (messageFilterObject: object(symKeyID: string, privateKeyID: string, sig: string, minPow: uint, topics: bytes4[]))" }
				}
				if (web3.shh.newSymKey) {
					web3.shh.newSymKey.toString = function() { return "function (callback)" }
				}
				if (web3.shh.post) {
					web3.shh.post.toString = function() { return "function (messageObject: object(symKeyId: string, pubKey: string, sig: string, ttl: uint, topic: string, payload: string, powTime: uint, powTarget: uint))" }
				}
				if (web3.shh.setMaxMessageSize) {
					web3.shh.setMaxMessageSize.toString = function() { return "function (size: uint, callback)" }
				}
				if (web3.shh.setMinPoW) {
					web3.shh.setMinPoW.toString = function() { return "function (pow: uint, callback)" }
				}
				if (web3.shh.version) {
					web3.shh.version.toString = function() { return "function ()" }
				}
			}
		}


		if (web3) {
			if (web3.admin) {
				if (web3.admin.addPeer) {
					web3.admin.addPeer.toString = function() { return "function (url(enode://<node-id>@<ip-address>:<port>): string)" }
				}
				if (web3.admin.addTrustedPeer) {
					web3.admin.addTrustedPeer.toString = function() { return "function (url(enode://<node-id>@<ip-address>:<port>): string)" }
				}
				if (web3.admin.exportChain) {
					web3.admin.exportChain.toString = function() { return "function (file: string, firstBlockNum: uint64, lastBlockNum: uint64)" }
				}
				if (web3.admin.importChain) {
					web3.admin.importChain.toString = function() { return "function (file: string)" }
				}
				if (web3.admin.removePeer) {
					web3.admin.removePeer.toString = function() { return "function (url(enode://<node-id>@<ip-address>:<port>): string)" }
				}
				if (web3.admin.removeTrustedPeer) {
					web3.admin.removeTrustedPeer.toString = function() { return "function (url(enode://<node-id>@<ip-address>:<port>): string)" }
				}
				if (web3.admin.startHTTP) {
					web3.admin.startHTTP.toString = function() { return "function (host: string, port: uint, cors: string, apis: string)" }
				}
				if (web3.admin.startRPC) {
					web3.admin.startRPC.toString = function() { return "function (host: string, port: uint, cors: string, apis: string)" }
				}
				if (web3.admin.startWS) {
					web3.admin.startWS.toString = function() { return "function (host: string, port: uint, cors: string, apis: string)" }
				}
			}
		}


		if (web3) {
			if (web3.bzz) {
				if (web3.bzz.download) {
					web3.bzz.download.toString = function() { return "function (bzzHash: string, localpath: string)" }
				}
				if (web3.bzz.get) {
					web3.bzz.get.toString = function() { return "function (bzzHash: string)" }
				}
				if (web3.bzz.modify) {
					web3.bzz.modify.toString = function() { return "function (manifestHash: string, path: string, data: stringOrBuffer, contentType: string)" }
				}
				if (web3.bzz.put) {
					web3.bzz.put.toString = function() { return "function (content: stringOrBuffer)" }
				}
				if (web3.bzz.store) {
					web3.bzz.store.toString = function() { return "function (data: stringOrBufferOrArrayBuffer)" }
				}
				if (web3.bzz.upload) {
					web3.bzz.upload.toString = function() { return "function (data: stringOrBuffer)" }
				}
			}
		}


		if (web3) {
			if (web3.BigNumber) {
				web3.BigNumber.toString = function() { return "function (number: int float string or BigNumber)" }
			}
			if (web3.createBatch) {
				web3.createBatch.toString = function() { return "function ()" }
			}
			if (web3.fromAscii) {
				web3.fromAscii.toString = function() { return "function (asciiData: string)" }
			}
			if (web3.fromDecimal) {
				web3.fromDecimal.toString = function() { return "function (value: int float or string)" }
			}
			if (web3.fromICAP) {
				web3.fromICAP.toString = function() { return "function (icap: string)" }
			}
			if (web3.fromUtf8) {
				web3.fromUtf8.toString = function() { return "function (utf8Data: string)" }
			}
			if (web3.fromWei) {
				web3.fromWei.toString = function() { return "function (number: uint string or BigNumber, unit: string)" }
			}
			if (web3.isAddress) {
				web3.isAddress.toString = function() { return "function (HexString: string)" }
			}
			if (web3.isChecksumAddress) {
				web3.isChecksumAddress.toString = function() { return "function (address: bytes)" }
			}
			if (web3.isConnected) {
				web3.isConnected.toString = function() { return "function ()" }
			}
			if (web3.padLeft) {
				web3.padLeft.toString = function() { return "function (value: string, length: uint, paddingCharacter: string)" }
			}
			if (web3.padRight) {
				web3.padRight.toString = function() { return "function (value: string, length: uint, paddingCharacter: string)" }
			}
			if (web3.reset) {
				web3.reset.toString = function() { return "function (keepIsSyncing: bool)" }
			}
			if (web3.setProvider) {
				web3.setProvider.toString = function() { return "function (provider: web3 provider)" }
			}
			if (web3.sha3) {
				web3.sha3.toString = function() { return "function (data: bytes)" }
			}
			if (web3.toAscii) {
				web3.toAscii.toString = function() { return "function (hex: string)" }
			}
			if (web3.toBigNumber) {
				web3.toBigNumber.toString = function() { return "function (value: int or string)" }
			}
			if (web3.toChecksumAddress) {
				web3.toChecksumAddress.toString = function() { return "function (address: bytes20)" }
			}
			if (web3.toDecimal) {
				web3.toDecimal.toString = function() { return "function (hexString: string)" }
			}
			if (web3.toHex) {
				web3.toHex.toString = function() { return "function (val: string int float object array or BigNumber)" }
			}
			if (web3.toUtf8) {
				web3.toUtf8.toString = function() { return "function (hexString: string)" }
			}
			if (web3.toWei) {
				web3.toWei.toString = function() { return "function (number: string or BigNumber, unit: string)" }
			}
		}
	`)
}
