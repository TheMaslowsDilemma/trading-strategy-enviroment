package simulation

import (
    "fmt"
    "time"
    "sync"
    "tse-p2/wallet"
    "tse-p2/token"
    "tse-p2/ledger"
)

type Simulation struct { 
    start       time.Time
    end         time.Time
    Dur         time.Duration
    CurrentDur  time.Duration
    CancelChan  chan byte
    LedgerLock  sync.Mutex
    Ledger      ledger.Ledger
}

func CreateSimulation(d time.Duration) (*Simulation, error) {
    return &Simulation{
        Dur: d,
        CurrentDur: 0,
        CancelChan: make(chan byte, 1),
        Ledger: make(ledger.Ledger),
    }, nil
}

func (s *Simulation) AddLedgerItem(id ledger.LedgerAddr, li ledger.LedgerItem) error {
    var existing ledger.LedgerItem

    s.LedgerLock.Lock()
    defer s.LedgerLock.Unlock()

    existing = s.Ledger[id]
    if existing != nil {
        return fmt.Errorf("ledger item already exists at %v", id)
    }

    s.Ledger[id] = li
    return nil
}

func (s *Simulation) AddWallet(initamnt uint64) ledger.LedgerAddr {
    var (
        walletAddr      ledger.LedgerAddr
        usdRsvAddr      ledger.LedgerAddr
        usdRsv          token.TokenReserve
        w               wallet.Wallet
        walletRsvs      []ledger.LedgerAddr
        err             error
    )
    
    for {
        // Add initial reserve to ledger
        usdRsvAddr = ledger.RandomLedgerAddr()
        usdRsv = token.TokenReserve {
            Amount: initamnt,
            Symbol: "usd",
        }
        err = s.AddLedgerItem(usdRsvAddr, usdRsv)
        if err != nil {
            continue
        }

        // Add the wallet to the ledger
        walletAddr = ledger.RandomLedgerAddr()
        walletRsvs = make([]ledger.LedgerAddr, 0)
        w = wallet.Wallet {
            TraderId: uint64(walletAddr),
            Reserves: walletRsvs,
        }
        w.AddReserve(usdRsvAddr)
        err = s.AddLedgerItem(walletAddr, w)
        if err != nil {
            continue
        }

        break
    }

    return walletAddr
}

func (s *Simulation) Run() {
    s.start = time.Now()

    for {
        select {
            case <-s.CancelChan:
                s.end = time.Now()
                s.CancelChan <- 0
                return
            default:
                s.Iter()
        }
    }
}

func (s *Simulation) Iter() {
    var cd time.Duration

    cd = time.Since(s.start)
    if cd >= s.Dur {
        s.CancelChan <- 1
        return
    } else {
        s.CurrentDur = cd
    }
}

func (s *Simulation) GetLedgerItemString(id ledger.LedgerAddr) (string, error) {
    var (
        str string
        err error
    )

    s.LedgerLock.Lock()
    str, err = s.Ledger.GetItemString(id)
    s.LedgerLock.Unlock()

    return str, err
}
