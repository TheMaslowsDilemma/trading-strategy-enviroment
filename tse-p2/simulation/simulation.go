package simulation

import (
    "fmt"
    "time"
    "sync"
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
        Ledger: make(map[uint64]ledger.LedgerItem),
    }, nil
}

func (s *Simulation) AddLedgerItem(id uint64, li ledger.LedgerItem) error {
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

func (s *Simulation) GetLedgerItemString(id uint64) (string, error) {
    var (
        str string
        err error
    )

    s.LedgerLock.Lock()
    str, err = s.Ledger.GetItemString(id)
    s.LedgerLock.Unlock()

    return str, err
}
