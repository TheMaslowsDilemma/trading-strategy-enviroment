package simulation

import (
	"fmt"
	"tse-p3/wallets"
	"tse-p3/exchanges"
	"tse-p3/ledger"
	"tse-p3/traders"
	"github.com/holiman/uint256"
	"github.com/cespare/xxhash"
)

type Simulation struct {
	MainLedger 		ledger.Ledger
	Traders			map[uint64] traders.Trader
	ExchangeDirectory	map[uint64] ledger.Addr
}

func NewSimulation() Simulation {
	return Simulation {
		MainLedger: ledger.CreateLedger(),
		Traders: make(map[uint64]traders.Trader),
		ExchangeDirectory: make(map[uint64]ledger.Addr),
	}
}	

func (s Simulation) String() string {
	return MainLedger.String()
}

func (s Simulation) GetNetworth(traderKey uint64) (*uint256.Int, error) {
	var (
		tr	trader.Trader
	)

	tr = s.Traders[traderKey]
	if tr.Id == 0 { 
		return nil, fmt.Errorf("no trader exists for key: %v", traderKey)
	}
	return tr.GetNetworth(s.rateProvider, s.walletProvider), nil
}
