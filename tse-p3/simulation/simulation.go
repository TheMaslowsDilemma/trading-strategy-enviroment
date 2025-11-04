package simulation

import (
	"fmt"
	"tse-p3/ledger"
	"tse-p3/traders"
	"tse-p3/users"
	"tse-p3/globals"
	"tse-p3/exchanges"
)

type Simulation struct {
	MainLedger 			ledger.Ledger
	Users				map[uint64] users.User
	Traders				map[uint64] traders.Trader
	ExchangeDirectory	map[uint64] ledger.Addr
}

func NewSimulation() Simulation {
	var (
		sim		Simulation
		cped	exchanges.CpeDescriptor
	)

	sim = Simulation {
		MainLedger: ledger.CreateLedger(),
		Users: make(map[uint64]users.User),
		Traders: make(map[uint64]traders.Trader),
		ExchangeDirectory: make(map[uint64]ledger.Addr),
	}

	// NOTE this adds the default exchange
	cped = exchanges.CpeDescriptor {
		AmountA: globals.TSECurrencyAmount,
		SymbolA: globals.TSECurrencySymbol,
		AmountB: globals.USDCurrencyAmount,
		SymbolB: globals.USDCurrencySymbol,
	}

	sim.MainLedger.AddConstantProductExchange(cped, 0) // NOTE 
	return sim
}

func (s Simulation) String() string {
	return fmt.Sprintf("{ ledger: %v; }", s.MainLedger)
}
