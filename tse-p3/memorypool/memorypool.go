package memorypool

import (
	"tse-p3/transactions"
)

type MemoryPool chan txs.Tx

func CreateMemoryPool(size uint) MemoryPool {
	return make(chan txs.Tx, size)
}

func (mp MemoryPool) Push(tx txs.Tx) bool {
	select {
	case mp <-tx:
		return true
	default:
		return false // NOTE we could tx.Notify(tx.FailedTx)...
	}
}

func (mp MemoryPool) Pop() (txs.Tx, bool) {
	var outtx txs.Tx
	select {
	case outtx = <-mp:
		return outtx, true
	default:
		return nil, false
	}
}


