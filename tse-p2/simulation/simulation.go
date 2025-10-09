package simulation

import (
    "time"
    "sync"
    "tse-p2/ledger"
    "tse-p2/mempool"
    "tse-p2/miner"
)

const simulationMemoryPoolSize = 512
const simulationEntityLogBufferSize = 10

type Simulation struct { 
    start       time.Time
    end         time.Time
    MaxDur      time.Duration
    RunningDur  time.Duration
    CancelChan  chan byte
    LedgerLock  sync.Mutex
    Ledger      ledger.Ledger
    MainMiner   miner.Miner
    MemoryPool  mempool.MemPool
}

func CreateSimulation(maxdur time.Duration) (*Simulation, error) {
    var (
        mm      miner.Miner
        lg      ledger.Ledger
        mp      mempool.MemPool
        cc      chan byte
    )

    lg = make(ledger.Ledger)
    cc = make(chan byte, 1)
    mm = miner.CreateMiner(lg, simulationEntityLogBufferSize)
    mp = mempool.CreateMempool(simulationMemoryPoolSize)

    return &Simulation{
        MaxDur: maxdur,
        RunningDur: 0,
        CancelChan: cc,
        Ledger: lg,
        MainMiner: mm,
        MemoryPool: mp,
    }, nil
}

func (s *Simulation) Run() {
    s.start = time.Now()

    // Start Entities
    go s.minerTask()


    // Enter Simulation Control Loop
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
    var currentd time.Duration

    currentd = time.Since(s.start)
    if currentd >= s.MaxDur {
        s.CancelChan <- 1
        return
    } else {
        s.RunningDur = currentd
    }
}

