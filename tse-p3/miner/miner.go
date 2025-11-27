package miner

import (
	//"fmt"
	"tse-p3/ledger"
	"tse-p3/globals"
	"tse-p3/transactions"
	"tse-p3/memorypool"
)


func NextBlock(tick uint64, mpl *memorypool.MemoryPool, scnd *ledger.Ledger) error {
	var (
		txblock	[]txs.Tx
		delta	*ledger.Ledger
		tx		txs.Tx
		popped	bool
		err		error
	)

	txblock = make([]txs.Tx, 0)
	for i := 0; i < globals.MaxBlockSize; i++ {
		tx, popped = mpl.Pop()
		if !popped {
			break // no more TXs in mem pool
		}
		txblock = append(txblock, tx)
	}

	for _, tx = range txblock {
		delta, err = tx.Apply(tick, *scnd)
		if err != nil {
			//fmt.Printf("Error applying tx: %v\n", err)
			tx.Notify(txs.TxFail)
			continue
		}
		tx.Notify(txs.TxPass)
		scnd.Merge(delta)
	}

	return  nil
}
