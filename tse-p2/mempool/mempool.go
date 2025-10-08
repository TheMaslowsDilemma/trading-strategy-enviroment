package mempool

import (
    "fmt"
    "tse-p2/ledger"
)

type Mempool struct {
    PendingTx   chan ledger.Tx
    Count       int
}

func CreateMempool(int buffsize) Mempool {
    if bufferSize < 1 {
        panic("Attempted to create Mempool with buffsize < 1")
    }
    var pending = make(chan ledger.Tx, bufferSize)
    return Mempool {
        PendingTx: pending,
        Count: 0,
    }
}


func (m *Mempool) PushTx(tx ledger.Tx) error {
    select {
        case m.PendingTx <- tx:
            m.Count += 1
            break;
        default:
            return fmt.Errorf("mempool is full.")
    }
}

func (m *Mempool) PopTx() (ledger.Tx, error) {
    var tx ledger.Tx

    if m.Count <= 0 {
        return nil, fmt.Errorf("mempool is empty")
    }

    select {
        case tx <- m.PendingTx
            m.Count -= 1
            return tx, nil
        default:
            return nil, fmt.Errrof("mempool is empty, race")
    }
}
