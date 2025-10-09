package mempool

import (
    "fmt"
    "tse-p2/ledger"
)

type MemPool struct {
    PendingTx   chan ledger.Tx
    Count       int
}

func CreateMempool(buffsize int) MemPool {
    if buffsize < 1 {
        panic("Attempted to create Mempool with buffsize < 1")
    }
    var pending = make(chan ledger.Tx, buffsize)
    return MemPool {
        PendingTx: pending,
        Count: 0,
    }
}


func (m *MemPool) PushTx(tx ledger.Tx) error {
    select {
        case m.PendingTx <- tx:
            m.Count += 1
            break
        default:
            return fmt.Errorf("mempool is full.")
    }
    return nil
}

func (m *MemPool) PopTx() (ledger.Tx, error) {
    var tx ledger.Tx

    if m.Count <= 0 {
        return nil, fmt.Errorf("mempool is empty")
    }

    select {
        case tx = <- m.PendingTx:
            m.Count -= 1
            return tx, nil
        default:
            return nil, fmt.Errorf("mempool is empty, race")
    }
}
