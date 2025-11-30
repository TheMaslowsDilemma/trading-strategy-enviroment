package simulation

import (
	"fmt"
	"sync"
	"tse-p3/ledger"
	"tse-p3/memorypool"
	"tse-p3/traders"
	"tse-p3/bots"
	"tse-p3/globals"
	"tse-p3/exchanges"
)

type Simulation struct {
	PrimaryLedger 		*ledger.Ledger
	SecondaryLedger		*ledger.Ledger
	PrimaryLock			sync.Mutex
	SecondaryLock		sync.Mutex
	TraderLock			sync.Mutex
	MemoryPool			memorypool.MemoryPool
	Bots				map[uint64] *bots.Bot
	Traders				map[uint64] *traders.Trader
	ExchangeDirectory	map[uint64] ledger.Addr
	CancelRequested		bool
}

func NewSimulation() Simulation {
	var (
		sim		Simulation
		cped	exchanges.CpeDescriptor
	)

	sim = Simulation {
		PrimaryLedger:		ledger.NewLedgerWithEmitter(),
		MemoryPool: 		memorypool.CreateMemoryPool(globals.DefaultMemoryPoolSize),
		Bots: 				make(map[uint64] *bots.Bot),
		Traders: 			make(map[uint64] *traders.Trader),
		ExchangeDirectory: 	make(map[uint64]ledger.Addr),
		CancelRequested: 	false,
	}

	// NOTE this adds the default exchange
	cped = exchanges.CpeDescriptor {
		AmountA: globals.TSECurrencyAmount,
		SymbolA: globals.TSESymbol,
		AmountB: globals.USDCurrencyAmount,
		SymbolB: globals.USDSymbol,
	}

	sim.AddExchange(cped, 0)

	cped = exchanges.CpeDescriptor {
		AmountA: globals.TSECurrencyAmount,
		SymbolA: "alpha",
		AmountB: globals.USDCurrencyAmount,
		SymbolB: globals.USDSymbol,
	}
	
	sim.AddExchange(cped, 0)

	cped = exchanges.CpeDescriptor {
		AmountA: globals.TSECurrencyAmount,
		SymbolA: "alpha",
		AmountB: globals.TSECurrencyAmount,
		SymbolB: globals.TSESymbol,
	}
	
	sim.AddExchange(cped, 0)
	sim.SecondaryLedger = ledger.Clone(*sim.PrimaryLedger)
	return sim
}

func (s Simulation) String() string {
	return fmt.Sprintf("{ ledger: %v; }", s.PrimaryLedger) 
}
