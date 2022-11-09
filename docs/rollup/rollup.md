## How To Rebuild A Chain From RollupInputBatch
As the whitepaper introduced, RollupInputBatch linked like this:

rollupInputBatch format in this:
```text
   // format: batchIndex(uint64)+ queueNum(uint64) + queueStartIndex(uint64) + subBatchNum(uint64) + subBatch0Time(uint64) +
    // subBatchLeftTimeDiff([]uint32) + batchesData
    // batchesData: version(0) + rlp([][]transaction)
```
```text
Batch_0 => Batch_1 => Batch_2 => Batch_3 ...
```

And the normal chain struct this:
```text
Block_0 => Block_1 => Block_2 => Block_3 =>Block_4 => Block_5 ...
```

and a rollupInputBatch can derive a list of blocks in this way:

#### first, sort transaction with context:
- transaction sorted in timestamp from low to high
- queue transaction sorted in queue index
- queue transaction always ahead of same timestamp l2 origin tx

### chose pending transactions for a block:

- a block can only have queue txs or origin l2 txs
- l2 origin tx must check first, the rule is same as validate tx.l1 origin tx do not need to check, which is already checked in l1.
- as origin miner, chose a list of txs to run for a new block. 
  - l2 origin tx with same block is pending txs for a block.
    - queue tx with same timestamp and ensure if all those queue txs succeed, will not run out of gas, then chosen queue txs is pending txs for a block.

### execute pending txs to generate a block:
same as evm logic, except:
- queue tx do not need nonce check, which is already checked in l1.
- queue tx's IntrinsicGas is zero.(these intrinstic gas is used for evaluate the cost to send tx to l1, but queue tx is already in l1)

### chain up this blocks
a block is next block's parent block. the first block that this batch generated parent is last batch's last block.the first batch's first block's parent is genesis block.