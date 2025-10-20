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
)

const simulationMemoryPoolSize = 512
const simulationEntityLogBufferSize = 10

type Simulation struct { 
    start       time.Time
    end         time.Time
    RunningDur  time.Duration
    IsCanceled  bool
    LedgerLock  sync.Mutex
    Ledger      ledger.Ledger
    MainMiner   miner.Miner
    MemoryPool  mempool.MemPool
    CliTrader   trader.Trader
    CliWallet   ledger.LedgerAddr
    ExAddr      ledger.LedgerAddr
    CandleNotifier func()
}

func CreateSimulation() (*Simulation, error) {
    var (
        mm      miner.Miner
        lg      ledger.Ledger
        mp      mempool.MemPool
        eaddr   ledger.LedgerAddr
        trdr    trader.Trader
        waddr   ledger.LedgerAddr
    )

    lg = make(ledger.Ledger)
    mp = mempool.CreateMempool(simulationMemoryPoolSize)
    eaddr = exchange.InitConstantProductExchange("usd", "eth", 10000, 500000, &lg)

    waddr = wallet.InitDefaultWallet(&lg)
    mm = miner.CreateMiner(simulationEntityLogBufferSize, lg)

    traderPlaceTx := func (tx ledger.Tx) {
        (&mp).PushTx(tx)
    }
    traderLedgerLookup := func (addr ledger.LedgerAddr) ledger.LedgerItem {
        return lg[addr]
    }
    trdr = trader.CreateTrader(
        nil, // no strategy for user trader
        10,
        waddr,
        eaddr,
        "usd",
        "eth",
        nil,
        traderLedgerLookup,
        traderPlaceTx,
    )

    sim := Simulation {
        RunningDur: 0,
        IsCanceled: false,
        Ledger: lg,
        MainMiner: mm,
        MemoryPool: mp,
        ExAddr: eaddr,
        CliTrader: trdr,
        CliWallet: waddr,
    }

    return &sim, nil
}

func (s *Simulation) Run() {
    s.start = time.Now()

    // Start Entity Routines
    s.initializeTraders()
    ledger.Merge(&s.MainMiner.BackLedger, s.Ledger)
    go s.minerTask()

    // Simulation Control Loop
    for {
        if s.IsCanceled {
            s.end = time.Now()
            return
        }
        s.Iter()
    }
}

func (s *Simulation) Iter() {
    var currentd time.Duration
    currentd = time.Since(s.start)
    s.RunningDur = currentd
}

