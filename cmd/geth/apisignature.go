package main

import "github.com/ethereum/go-ethereum/console"

func injectApiSignature(c *console.Console) {
	c.Evaluate(`
		if (typeof l2 !== 'undefined') {  
			if (typeof l2.getBatch !== 'undefined') {  
				l2.getBatch.toString = function() { return "function (batchIndex: uint64, useDetail: bool)" }
			}
		}
	`)
}
