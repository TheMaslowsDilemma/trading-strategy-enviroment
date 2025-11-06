package simulation

import (
	"fmt"
	"sync"
	"tse-p3/ledger"
	"tse-p3/memorypool"
	"tse-p3/traders"
	"tse-p3/users"
	"tse-p3/globals"
	"tse-p3/exchanges"
	"tse-p3/miner"
)

type Simulation struct {
	MainLedger 			ledger.Ledger
	ScndLedger			ledger.Ledger
	LedgerLock			sync.Mutex
	MemoryPool			memorypool.MemoryPool
	Users				map[uint64] users.User
	Traders				map[uint64] traders.Trader
	ExchangeDirectory	map[uint64] ledger.Addr
	CancelRequested		bool
}

func NewSimulation() Simulation {
	var (
		sim		Simulation
		cped	exchanges.CpeDescriptor
	)

	sim = Simulation {
		MainLedger: ledger.CreateLedger(),
		MemoryPool: memorypool.CreateMemoryPool(globals.DefaultMemoryPoolSize),
		Users: make(map[uint64]users.User),
		Traders: make(map[uint64]traders.Trader),
		ExchangeDirectory: make(map[uint64]ledger.Addr),
		CancelRequested: false,
	}

	// NOTE this adds the default exchange
	cped = exchanges.CpeDescriptor {
		AmountA: globals.TSECurrencyAmount,
		SymbolA: globals.TSECurrencySymbol,
		AmountB: globals.USDCurrencyAmount,
		SymbolB: globals.USDCurrencySymbol,
	}
	sim.addExchange(cped, 0)
	sim.ScndLedger = miner.CreateSecondary(sim.MainLedger)
	return sim
}

func (s Simulation) String() string {
	return fmt.Sprintf("{ ledger: %v; }", s.MainLedger)
}
