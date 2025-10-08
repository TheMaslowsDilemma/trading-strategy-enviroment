package miner

import (
    "fmt"
    "tse-p2/ledger"
)

const MaxBlockSize = 7

type Miner struct { 
    BackLedger  ledger.Ledger
    Logs        chan string
}

func CreateMiner(lg ledger.Ledger, logsize int) Miner {
    var (
        bckldg  ledger.Ledger
        logs  chan string
        err     error
    )

    bckldg = make(ledger.Ledger)
    logs = make(chan string, logsize)

    err = ledger.MergeUpdate(bckldg, lg)
    if err != nil {
        panic(fmt.Error("Miner constructor bad ledger copy: %v", err))
    }
    
    return Miner {
        BackLedger: bckldg,
        Logs: logs,
    }
}

func (m *Miner) PushLog(log string) {
    select {
        case m.Logs <- log:
            break;
        default:
            <- m.Logs
            m.Logs <- log
            break;
    }
}

func (m *Miner) MineNextBlock(mpl mempool.Mempool) error {
    var (
        txs     []ledger.Tx
        tx      ledger.Tx
        lgp     ledger.Ledger
        err     error
    )

    // create txs
    txs, err = createTxs(mpl)
    if err != nil {
        return nil, fmt.Errorf("failed create txs: %v", err)
    }

    for _, tx = range txs {
        lgp, err = tx.Apply(m.BackLedger)
        if err != nil {
            m.PushLog(fmt.Sprintf("tx apply failed, skipping: %v", err))
            continue
        }
        ledger.MergeUpdate(m.BackLedger, lgp)
    }
}

func createTxBlock(mpl *mempool.MemoryPool) ([]ledger.Tx, error) {
    var (
        i       int
        tx      ledger.Tx
        txs     []ledger.Tx 
        err     error
    )               
    
    if mpl.Count == 0 {
        return [], fmt.Errorf("no tx to process")
    }

    txs = make([]ledger.Tx, 0)

    i = 0
    while i < MaxBlockSize {
        tx, err = mpl.PopTx()
        if err != nil {
            break // TODO Debug logging ?
        }
        txs = append(txs, x)
    }

    return txs, nil
}
