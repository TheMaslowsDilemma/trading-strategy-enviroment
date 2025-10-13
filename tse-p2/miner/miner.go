package miner

import (
    "fmt"
    "time"
    "tse-p2/ledger"
    "tse-p2/mempool"
)

const MaxBlockSize = 7

type Miner struct {
    LastBlockTime       time.Time
    BackLedger          ledger.Ledger
    Logs                chan string
}

func CreateMiner(logsize int, lg ledger.Ledger) Miner {
    var (
        bckldg  ledger.Ledger
        logs    chan string
    )

    bckldg = make(ledger.Ledger)
    logs   = make(chan string, logsize)
    ledger.Merge(&bckldg, lg)
    
    return Miner {
        LastBlockTime: time.Now(),
        BackLedger: bckldg,
        Logs: logs,
    }
}

func (m *Miner) PushLog(log string) {
    select {
        case m.Logs <- log:
            break
        default:
            <- m.Logs
            m.Logs <- log
            break
    }
}

func (m *Miner) MineNextBlock(tick uint64, mpl *mempool.MemPool) error {
    var (
        txs     []ledger.Tx
        tx      ledger.Tx
        lgp     ledger.Ledger
        err     error
    )

    txs, err = createTxBlock(mpl)
    if err != nil {
        return fmt.Errorf("failed create txs: %v", err)
    }

    for _, tx = range txs {
        lgp, err = tx.Apply(tick, m.BackLedger)
        if err != nil {
            fmt.Printf("tx apply failed, skipping: %v", err)
            continue
        }
        ledger.Merge(&m.BackLedger, lgp)
    }

    return nil
}

func createTxBlock(mpl *mempool.MemPool) ([]ledger.Tx, error) {
    var (
        i       int
        tx      ledger.Tx
        txs     []ledger.Tx 
        err     error
    )               
    
    if mpl.Count == 0 {
        return []ledger.Tx{}, fmt.Errorf("no tx to process")
    }

    txs = make([]ledger.Tx, 0)
    i = 0 
    for i < MaxBlockSize {
        
        tx, err = mpl.PopTx()
        if err != nil {
            break // TODO dbg log
        }
        txs = append(txs, tx)
        i += 1
    }

    return txs, nil
}
