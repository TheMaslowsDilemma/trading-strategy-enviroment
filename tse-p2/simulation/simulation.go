package simulation

import (
    "time"
    "sync"
    "tse-p2/ledger"
    "tse-p2/mempool"
    "tse-p2/miner"
    "tse-p2/exchange"
    "tse-p2/trader"
    "tse-p2/wallet"
    "tse-p2/token"
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
    CliTrader   trader.Trader
    CliWallet   ledger.LedgerAddr
    ExAddr      ledger.LedgerAddr
}

func CreateSimulation(maxdur time.Duration) (*Simulation, error) {
    var (
        mm      miner.Miner
        lg      ledger.Ledger
        mp      mempool.MemPool
        eaddr   ledger.LedgerAddr
        cc      chan byte
        trdr    trader.Trader
        waddr   ledger.LedgerAddr
    )

    lg = make(ledger.Ledger)
    cc = make(chan byte, 1)
    
    mm = miner.CreateMiner(lg, simulationEntityLogBufferSize)
    mp = mempool.CreateMempool(simulationMemoryPoolSize)
    eaddr = exchange.InitConstantProductExchange("usd", "eth", 2000, 500000000, lg)
    
    // Initialize CLI User Wallet and Trader //
    rs := []token.TokenReserve {
        token.TokenReserve {
            Symbol: "usd",
            Amount: 10000.0,
        },
        token.TokenReserve {
            Symbol: "eth",
            Amount: 0.0,
        },
    }

    waddr = wallet.InitWallet(rs, lg)
    trdr = trader.CreateTrader(
        nil, // no strategy for user trader
        10,
        waddr,
        eaddr,
        "usd",
        "eth",
        lg,
    )

    return &Simulation{
        MaxDur: maxdur,
        RunningDur: 0,
        CancelChan: cc,
        Ledger: lg,
        MainMiner: mm,
        MemoryPool: mp,
        ExAddr: eaddr,
        CliTrader: trdr,
        CliWallet: waddr,
    }, nil
}

func (s *Simulation) Run() {
    s.start = time.Now()

    // Start Entity Routines
    go s.minerTask()

    // Simulation Control Loop
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

