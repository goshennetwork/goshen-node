package main

import "github.com/ethereum/go-ethereum/console"

func injectApiSignature(c *console.Console) {
	c.Evaluate(`
		if (typeof l2 !== 'undefined') {  
			if (typeof l2.getBatch !== 'undefined') {  
				l2.getBatch.toString = function() { return "function (batchIndex: uint64, useDetail: bool)" }
			}
			if (typeof l2.getBatchIndex !== 'undefined') {  
				l2.getBatchIndex.toString = function() { return "function (blockNumber: uint64)" }
			}
			if (typeof l2.getBatchState !== 'undefined') {  
				l2.getBatchState.toString = function() { return "function (batchIndex: uint64)" }
			}
			if(typeof l2.getEnqueuedTxs !== 'undefined'){
				l2.getEnqueuedTxs.toString = function() { return "function (queueStart: uint64, queueNum: uint64)" }
			}
			if(typeof l2.getL2MMRProof !== 'undefined'){
				l2.getL2MMRProof.toString = function() { return "function (msgIndex: uint64, size: uint64)" }
			}
			if(typeof l2.getL2RelayMsgParams !== 'undefined'){
				l2.getL2RelayMsgParams.toString = function() { return "function (msgIndex: uint64)" }
			}
			if(typeof l2.getPendingTxBatches !== 'undefined'){
				l2.getPendingTxBatches.toString = function() { return "function ()" }
			}
			if(typeof l2.globalInfo !== 'undefined'){
				l2.globalInfo.toString = function() { return "function ()" }
			}
			if(typeof l2.getState !== 'undefined'){
				l2.getState.toString = function() { return "function (batchIndex: uint64)" }
			}
			if(typeof l2.inputBatchNumber !== 'undefined'){
				l2.inputBatchNumber.toString = function() { return "function ()" }
			}
			if(typeof l2.stateBatchNumber !== 'undefined'){
				l2.stateBatchNumber.toString = function() { return "function ()" }
			}
		}


		if(typeof eth !== 'undefined'){
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


		if(typeof debug !== 'undefined'){
			if (typeof debug.accountRange !== 'undefined') {
				debug.accountRange.toString = function() { return "function (blockNumberOrBlockHash: uint latest earliest pending or string, start: string, maxResults: uint, skip: uint, reverse: bool)" }
			}
			if (typeof debug.backtraceAt !== 'undefined') {
				debug.backtraceAt.toString = function() { return "function ((file:line): string)" }
			}
			if (typeof debug.blockProfile !== 'undefined') {
				debug.blockProfile.toString = function() { return "function (file: string, time: uint)" }
			}
			if (typeof debug.chaindbCompact !== 'undefined') {
				debug.chaindbCompact.toString = function() { return "function ()" }
			}
			if (typeof debug.chaindbProperty !== 'undefined') {
				debug.chaindbProperty.toString = function() { return "function (leveldb properties: string)" }
			}
			if (typeof debug.cpuProfile !== 'undefined') {
				debug.cpuProfile.toString = function() { return "function (file: string, time: uint)" }
			}
			if (typeof debug.dumpBlock !== 'undefined') {
				debug.dumpBlock.toString = function() { return "function (blockNumber: uint or latest earliest pending)" }
			}
			if (typeof debug.freeOSMemory !== 'undefined') {
				debug.freeOSMemory.toString = function() { return "function ()" }
			}
			if (typeof debug.freezeClient !== 'undefined') {
				debug.freezeClient.toString = function() { return "function ()" }
			}
			if (typeof debug.gcStats !== 'undefined') {
				debug.gcStats.toString = function() { return "function ()" }
			}
			if (typeof debug.getAccessibleState !== 'undefined') {
				debug.getAccessibleState.toString = function() { return "function (from: uint, to: uint)" }
			}
			if (typeof debug.getBadBlocks !== 'undefined') {
				debug.getBadBlocks.toString = function() { return "function ()" }
			}
			if (typeof debug.getBlockRlp !== 'undefined') {
				debug.getBlockRlp.toString = function() { return "function (blockNumberOrBlockHash: uint latest earliest pending or string)" }
			}
			if (typeof debug.getHeaderRlp !== 'undefined') {
				debug.getHeaderRlp.toString = function() { return "function (blockNumber: uint)" }
			}
			if (typeof debug.getModifiedAccountsByHash !== 'undefined') {
				debug.getModifiedAccountsByHash.toString = function() { return "function (startHash: string, endHash: string)" }
			}
			if (typeof debug.getModifiedAccountsByNumber !== 'undefined') {
				debug.getModifiedAccountsByNumber.toString = function() { return "function (startNum: uint64, endNum: uint64)" }
			}
			if (typeof debug.goTrace !== 'undefined') {
				debug.goTrace.toString = function() { return "function (file: string, runningTime: uint)" }
			}
			if (typeof debug.intermediateRoots !== 'undefined') {
				debug.intermediateRoots.toString = function() { return "function (blockHash: string)" }
			}
			if (typeof debug.mutexProfile !== 'undefined') {
				debug.mutexProfile.toString = function() { return "function (file: string, nsec: uint)" }
			}
			if (typeof debug.preimage !== 'undefined') {
				debug.preimage.toString = function() { return "function (hash: string)" }
			}
			if (typeof debug.printBlock !== 'undefined') {
				debug.printBlock.toString = function() { return "function (blockNumber: uint)" }
			}
			if (typeof debug.seedHash !== 'undefined') {
				debug.seedHash.toString = function() { return "function (blockNumber: uint)" }
			}
			if (typeof debug.setBlockProfileRate !== 'undefined') {
				debug.setBlockProfileRate.toString = function() { return "function (rate: int)" }
			}
			if (typeof debug.setGCPercent !== 'undefined') {
				debug.setGCPercent.toString = function() { return "function (v: int)" }
			}
			if (typeof debug.setHead !== 'undefined') {
				debug.setHead.toString = function() { return "function (blockNumber: string)" }
			}
			if (typeof debug.setMutexProfileFraction !== 'undefined') {
				debug.setMutexProfileFraction.toString = function() { return "function (rate: int)" }
			}
			if (typeof debug.stacks !== 'undefined') {
				debug.stacks.toString = function() { return "function ()" }
			}
			if (typeof debug.standardTraceBadBlockToFile !== 'undefined') {
				debug.standardTraceBadBlockToFile.toString = function() { return "function (blockHash: string, optionObject: object(disableStorage: bool, disableMemory: bool, disableStack: bool, fullStorage: bool))" }
			}
			if (typeof debug.standardTraceBlockToFile !== 'undefined') {
				debug.standardTraceBlockToFile.toString = function() { return "function (blockHash: string, optionObject: object(disableStorage: bool, disableMemory: bool, disableStack: bool, fullStorage: bool))" }
			}
			if (typeof debug.startCPUProfile !== 'undefined') {
				debug.startCPUProfile.toString = function() { return "function (file: string)" }
			}
			if (typeof debug.startGoTrace !== 'undefined') {
				debug.startGoTrace.toString = function() { return "function (file: string)" }
			}
			if (typeof debug.stopCPUProfile !== 'undefined') {
				debug.stopCPUProfile.toString = function() { return "function ()" }
			}
			if (typeof debug.stopGoTrace !== 'undefined') {
				debug.stopGoTrace.toString = function() { return "function ()" }
			}
			if (typeof debug.storageRangeAt !== 'undefined') {
				debug.storageRangeAt.toString = function() { return "function (blockHash: string, transactionIndex: uint, contractAddress: bytes20, keyStart: string, maxResult: uint)" }
			}
			if (typeof debug.testSignCliqueBlock !== 'undefined') {
				debug.testSignCliqueBlock.toString = function() { return "function (sigAddress: bytes20, blockNumber: uint)" }
			}
			if (typeof debug.traceBadBlock !== 'undefined') {
				debug.traceBadBlock.toString = function() { return "function (blockHash: string, optionObject: object(disableStorage: bool, disableStack: bool, enableMemory: bool, enableReturnData: bool, tracer: string, timeOut: string, tracerConfig: object(onlyTopCall: bool, withLog: bool)))" }
			}
			if (typeof debug.traceBlock !== 'undefined') {
				debug.traceBlock.toString = function() { return "function (blockRlp: string, optionObject: object(disableStorage: bool, disableStack: bool, enableMemory: bool, enableReturnData: bool, tracer: string, timeOut: string, tracerConfig: object(onlyTopCall: bool, withLog: bool)))" }
			}
			if (typeof debug.traceBlockByHash !== 'undefined') {
				debug.traceBlockByHash.toString = function() { return "function (blockHash: string, optionObject: object(disableStorage: bool, disableStack: bool, enableMemory: bool, enableReturnData: bool, tracer: string, timeOut: string, tracerConfig: object(onlyTopCall: bool, withLog: bool)))" }
			}
			if (typeof debug.traceBlockByNumber !== 'undefined') {
				debug.traceBlockByNumber.toString = function() { return "function (blockNumber: uint, optionObject: object(disableStorage: bool, disableStack: bool, enableMemory: bool, enableReturnData: bool, tracer: string, timeOut: string, tracerConfig: object(onlyTopCall: bool, withLog: bool)))" }
			}
			if (typeof debug.traceBlockFromFile !== 'undefined') {
				debug.traceBlockFromFile.toString = function() { return "function (file: string, optionObject: object(disableStorage: bool, disableStack: bool, enableMemory: bool, enableReturnData: bool, tracer: string, timeOut: string, tracerConfig: object(onlyTopCall: bool, withLog: bool)))" }
			}
			if (typeof debug.traceCall !== 'undefined') {
				debug.traceCall.toString = function() { return "function ((transaction: object(from: bytes20, to: bytes20, gas: uint, gasPrice: uint, value: uint, data: bytes, nonce: uint), blockNumberOrBlockHash: uint latest earliest pending or string, optionObject: (disableStorage: bool, disableStack: bool, enableMemory: bool, enableReturnData: bool, tracer: string, timeOut: string, tracerConfig: object(onlyTopCall: bool, withLog: bool)))" }
			}
			if (typeof debug.traceTransaction !== 'undefined') {
				debug.traceTransaction.toString = function() { return "function (transactionHash: string, optionObject: object(disableStorage: bool, disableStack: bool, enableMemory: bool, enableReturnData: bool, tracer: string, timeOut: string, tracerConfig: object(onlyTopCall: bool, withLog: bool)))" }
			}
			if (typeof debug.verbosity !== 'undefined') {
				debug.verbosity.toString = function() { return "function (logLevel: int)" }
			}
			if (typeof debug.vmodule !== 'undefined') {
				debug.vmodule.toString = function() { return "function (logPattern: string)" }
			}
			if (typeof debug.writeBlockProfile !== 'undefined') {
				debug.writeBlockProfile.toString = function() { return "function (file: string, seconds: uint)" }
			}
			if (typeof debug.writeMemProfile !== 'undefined') {
				debug.writeMemProfile.toString = function() { return "function (file: string)" }
			}
			if (typeof debug.writeMutexProfile !== 'undefined') {
				debug.writeMutexProfile.toString = function() { return "function (file: string)" }
			}
			if (typeof debug.getReadStorageProofAtBlock !== 'undefined') {
				debug.getReadStorageProofAtBlock.toString = function() { return "function (blockNumber: uint64)" }
			}
		}


		if(typeof miner !== 'undefined'){
			if(typeof miner.getHashrate !== 'undefined'){
				miner.getHashrate.toString = function() { return "function ()" }
			}
			if (typeof miner.setEtherbase !== 'undefined') {
				miner.setEtherbase.toString = function() { return "function (address: bytes20)" }
			}
			if (typeof miner.setExtra !== 'undefined') {
				miner.setExtra.toString = function() { return "function (extraData: bytes32)" }
			}
			if (typeof miner.setGasLimit !== 'undefined') {
				miner.setGasLimit.toString = function() { return "function (gasLimit: uint)" }
			}
			if (typeof miner.setGasPrice !== 'undefined') {
				miner.setGasPrice.toString = function() { return "function (gasPrice: uint)" }
			}
			if (typeof miner.setRecommitInterval !== 'undefined') {
				miner.setRecommitInterval.toString = function() { return "function (timeInterval: uint)" }
			}
			if (typeof miner.start !== 'undefined') {
				miner.start.toString = function() { return "function (threadsNum: unit)" }
			}
			if (typeof miner.stop !== 'undefined') {
				miner.stop.toString = function() { return "function ()" }
			}
		}


		if(typeof personal !== 'undefined'){
			if (typeof personal.deriveAccount !== 'undefined') {
				personal.deriveAccount.toString = function() { return "function (seed: string, derivationPath: string)" }
			}
			if (typeof personal.ecRecover !== 'undefined') {
				personal.ecRecover.toString = function() { return "function (data: string, signature: string)" }
			}
			if (typeof personal.getListAccounts !== 'undefined') {
				personal.getListAccounts.toString = function() { return "function (callback)" }
			}
			if (typeof personal.getListWallets !== 'undefined') {
				personal.getListWallets.toString = function() { return "function (callback)" }
			}
			if (typeof personal.importRawKey !== 'undefined') {
				personal.importRawKey.toString = function() { return "function (privateKey: string, passphrase: string)" }
			}
			if (typeof personal.initializeWallet !== 'undefined') {
				personal.initializeWallet.toString = function() { return "function (url: string)" }
			}
			if (typeof personal.lockAccount !== 'undefined') {
				personal.lockAccount.toString = function() { return "function (address: bytes20)" }
			}
			if (typeof personal.newAccount !== 'undefined') {
				personal.newAccount.toString = function() { return "function (password: string)" }
			}
			if (typeof personal.openWallet !== 'undefined') {
				personal.openWallet.toString = function() { return "function (url: string, passphrase: string)" }
			}
			if (typeof personal.sendTransaction !== 'undefined') {
				personal.sendTransaction.toString = function() { return "function (transaction: object(from: bytes20, to: bytes20, gas: uint, gasPrice: uint, value: uint, data: bytes, nonce: uint), passphrase: string)" }
			}
			if (typeof personal.sign !== 'undefined') {
				personal.sign.toString = function() { return "function (dataToSign: string, accountAddress: bytes20, passphrase: string)" }
			}
			if (typeof personal.signTransaction !== 'undefined') {
				personal.signTransaction.toString = function() { return "function (transaction: object(from: bytes20, to: bytes20, gas: uint, gasPrice: uint, value: uint, data: bytes, nonce: uint), passphrase: string)" }
			}
			if (typeof personal.unlockAccount !== 'undefined') {
				personal.unlockAccount.toString = function() { return "function (accountAddress: bytes20, passphrase: string, time: uint)" }
			}
			if (typeof personal.unpair !== 'undefined') {
				personal.unpair.toString = function() { return "function (url: string, pin: string)" }
			}
		}


		if(typeof utils !== 'undefined'){
			if (typeof utils.saveString !== 'undefined') {
				utils.saveString.toString = function() { return "function (data: string, fileName: string)" }
			}
		}


		if(typeof web3 !== 'undefined'){
			if (typeof web3.shh !== 'undefined') {
				if (typeof web3.shh.addPrivateKey !== 'undefined') {
					web3.shh.addPrivateKey.toString = function() { return "function (privateKey: string, callback)" }
				}
				if (typeof web3.shh.addSymKey !== 'undefined') {
					web3.shh.addSymKey.toString = function() { return "function (symKey: string, callback)" }
				}
				if (typeof web3.shh.deleteKeyPair !== 'undefined') {
					web3.shh.deleteKeyPair.toString = function() { return "function (id: string, callback)" }
				}
				if (typeof web3.shh.deleteSymKey !== 'undefined') {
					web3.shh.deleteSymKey.toString = function() { return "function (symKeyID: string, callback)" }
				}
				if (typeof web3.shh.generateSymKeyFromPassword !== 'undefined') {
					web3.shh.generateSymKeyFromPassword.toString = function() { return "function (password: string, callback)" }
				}
				if (typeof web3.shh.getPrivateKey !== 'undefined') {
					web3.shh.getPrivateKey.toString = function() { return "function (id: string, callback)" }
				}
				if (typeof web3.shh.getPublicKey !== 'undefined') {
					web3.shh.getPublicKey.toString = function() { return "function (id: string, callback)" }
				}
				if (typeof web3.shh.getSymKey !== 'undefined') {
					web3.shh.getSymKey.toString = function() { return "function (id: string, callback)" }
				}
				if (typeof web3.shh.hasKeyPair !== 'undefined') {
					web3.shh.hasKeyPair.toString = function() { return "function (id: string, callback)" }
				}
				if (typeof web3.shh.hasSymKey !== 'undefined') {
					web3.shh.hasSymKey.toString = function() { return "function (id: string, callback)" }
				}
				if (typeof web3.shh.markTrustedPeer !== 'undefined') {
					web3.shh.markTrustedPeer.toString = function() { return "function (enode: string, callback)" }
				}
				if (typeof web3.shh.newKeyPair !== 'undefined') {
					web3.shh.newKeyPair.toString = function() { return "function (callback)" }
				}
				if (typeof web3.shh.newMessageFilter !== 'undefined') {
					web3.shh.newMessageFilter.toString = function() { return "function (messageFilterObject: object(symKeyID: string, privateKeyID: string, sig: string, minPow: uint, topics: bytes4[]))" }
				}
				if (typeof web3.shh.newSymKey !== 'undefined') {
					web3.shh.newSymKey.toString = function() { return "function (callback)" }
				}
				if (typeof web3.shh.post !== 'undefined') {
					web3.shh.post.toString = function() { return "function (messageObject: object(symKeyId: string, pubKey: string, sig: string, ttl: uint, topic: string, payload: string, powTime: uint, powTarget: uint))" }
				}
				if (typeof web3.shh.setMaxMessageSize !== 'undefined') {
					web3.shh.setMaxMessageSize.toString = function() { return "function (size: uint, callback)" }
				}
				if (typeof web3.shh.setMinPoW !== 'undefined') {
					web3.shh.setMinPoW.toString = function() { return "function (pow: uint, callback)" }
				}
				if (typeof web3.shh.version !== 'undefined') {
					web3.shh.version.toString = function() { return "function ()" }
				}
			}
		}


		if(typeof web3 !== 'undefined'){
			if(typeof web3.admin !== 'undefined'){
				if(typeof web3.admin.addPeer !== 'undefined'){
					web3.admin.addPeer.toString = function() { return "function (url(enode://<node-id>@<ip-address>:<port>): string)" }
				}
				if (typeof web3.admin.addTrustedPeer !== 'undefined') {
					web3.admin.addTrustedPeer.toString = function() { return "function (url(enode://<node-id>@<ip-address>:<port>): string)" }
				}
				if (typeof web3.admin.exportChain !== 'undefined') {
					web3.admin.exportChain.toString = function() { return "function (file: string, firstBlockNum: uint64, lastBlockNum: uint64)" }
				}
				if (typeof web3.admin.importChain !== 'undefined') {
					web3.admin.importChain.toString = function() { return "function (file: string)" }
				}
				if (typeof web3.admin.removePeer !== 'undefined') {
					web3.admin.removePeer.toString = function() { return "function (url(enode://<node-id>@<ip-address>:<port>): string)" }
				}
				if (typeof web3.admin.removeTrustedPeer !== 'undefined') {
					web3.admin.removeTrustedPeer.toString = function() { return "function (url(enode://<node-id>@<ip-address>:<port>): string)" }
				}
				if (typeof web3.admin.startHTTP !== 'undefined') {
					web3.admin.startHTTP.toString = function() { return "function (host: string, port: uint, cors: string, apis: string)" }
				}
				if (typeof web3.admin.startRPC !== 'undefined') {
					web3.admin.startRPC.toString = function() { return "function (host: string, port: uint, cors: string, apis: string)" }
				}
				if (typeof web3.admin.startWS !== 'undefined') {
					web3.admin.startWS.toString = function() { return "function (host: string, port: uint, cors: string, apis: string)" }
				}
			}
		}

		if(typeof web3 !== 'undefined'){
			if(typeof web3.bzz !== 'undefined'){
				if (typeof web3.bzz.download !== 'undefined') {
					web3.bzz.download.toString = function() { return "function (bzzHash: string, localpath: string)" }
				}
				if (typeof web3.bzz.get !== 'undefined') {
					web3.bzz.get.toString = function() { return "function (bzzHash: string)" }
				}
				if (typeof web3.bzz.modify !== 'undefined') {
					web3.bzz.modify.toString = function() { return "function (manifestHash: string, path: string, data: stringOrBuffer, contentType: string)" }
				}
				if (typeof web3.bzz.put !== 'undefined') {
					web3.bzz.put.toString = function() { return "function (content: stringOrBuffer)" }
				}
				if (typeof web3.bzz.store !== 'undefined') {
					web3.bzz.store.toString = function() { return "function (data: stringOrBufferOrArrayBuffer)" }
				}
				if (typeof web3.bzz.upload !== 'undefined') {
					web3.bzz.upload.toString = function() { return "function (data: stringOrBuffer)" }
				}
			}
		}


		if(typeof web3 !== 'undefined'){
			if (typeof web3.BigNumber !== 'undefined') {
				web3.BigNumber.toString = function() { return "function (number: int float string or BigNumber)" }
			}
			if (typeof web3.createBatch !== 'undefined') {
				web3.createBatch.toString = function() { return "function ()" }
			}
			if (typeof web3.fromAscii !== 'undefined') {
				web3.fromAscii.toString = function() { return "function (asciiData: string)" }
			}
			if (typeof web3.fromDecimal !== 'undefined') {
				web3.fromDecimal.toString = function() { return "function (value: int float or string)" }
			}
			if (typeof web3.fromICAP !== 'undefined') {
				web3.fromICAP.toString = function() { return "function (icap: string)" }
			}
			if (typeof web3.fromUtf8 !== 'undefined') {
				web3.fromUtf8.toString = function() { return "function (utf8Data: string)" }
			}
			if (typeof web3.fromWei !== 'undefined') {
				web3.fromWei.toString = function() { return "function (number: uint string or BigNumber, unit: string)" }
			}
			if (typeof web3.isAddress !== 'undefined') {
				web3.isAddress.toString = function() { return "function (HexString: string)" }
			}
			if (typeof web3.isChecksumAddress !== 'undefined') {
				web3.isChecksumAddress.toString = function() { return "function (address: bytes)" }
			}
			if (typeof web3.isConnected !== 'undefined') {
				web3.isConnected.toString = function() { return "function ()" }
			}
			if (typeof web3.padLeft !== 'undefined') {
				web3.padLeft.toString = function() { return "function (value: string, length: uint, paddingCharacter: string)" }
			}
			if (typeof web3.padRight !== 'undefined') {
				web3.padRight.toString = function() { return "function (value: string, length: uint, paddingCharacter: string)" }
			}
			if (typeof web3.reset !== 'undefined') {
				web3.reset.toString = function() { return "function (keepIsSyncing: bool)" }
			}
			if (typeof web3.setProvider !== 'undefined') {
				web3.setProvider.toString = function() { return "function (provider: web3 provider)" }
			}
			if (typeof web3.sha3 !== 'undefined') {
				web3.sha3.toString = function() { return "function (data: bytes)" }
			}
			if (typeof web3.toAscii !== 'undefined') {
				web3.toAscii.toString = function() { return "function (hex: string)" }
			}
			if (typeof web3.toBigNumber !== 'undefined') {
				web3.toBigNumber.toString = function() { return "function (value: int or string)" }
			}
			if (typeof web3.toChecksumAddress !== 'undefined') {
				web3.toChecksumAddress.toString = function() { return "function (address: bytes20)" }
			}
			if (typeof web3.toDecimal !== 'undefined') {
				web3.toDecimal.toString = function() { return "function (hexString: string)" }
			}
			if (typeof web3.toHex !== 'undefined') {
				web3.toHex.toString = function() { return "function (val: string int float object array or BigNumber)" }
			}
			if (typeof web3.toUtf8 !== 'undefined') {
				web3.toUtf8.toString = function() { return "function (hexString: string)" }
			}
			if (typeof web3.toWei !== 'undefined') {
				web3.toWei.toString = function() { return "function (number: string or BigNumber, unit: string)" }
			}
		}	
	`)
}
