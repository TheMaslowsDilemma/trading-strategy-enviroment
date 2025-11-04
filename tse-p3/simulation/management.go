package simulation

import (
	"tse-p3/traders"
	"tse-p3/exchanges"
	"tse-p3/ledger"
	"tse-p3/wallets"
	"tse-p3/globals"
)

func (s *Simulation) AddUser() uint64 {
	// create trader
	// create wallet --> tse default + default amount

	var (
		trdr	traders.Trader
		wd		wallets.WalletDescriptor
		waddr	ledger.Addr
	)

	trdr = traders.CreateTrader()
	wd = wallets.WalletDescriptor {
		Amount: globals.UserStartingBalance,
		Symbol: globals.TSECurrencySymbol,
	}

	waddr = s.AddWallet(wd)
	trdr.AddWallet(wd.Symbol, waddr)
	s.AddTrader(trdr)

	return trdr.Id
}

func (s *Simulation) AddBot() {
	// TODO
}

func (s *Simulation) AddTrader(t traders.Trader) {
	s.Traders[t.Id] = t
}

func (s *Simulation) AddWallet(wd wallets.WalletDescriptor) ledger.Addr {
	return s.MainLedger.AddWallet(wd)
}

func (s *Simulation) AddExchange(cped exchanges.CpeDescriptor, tick uint64) ledger.Addr {
	return s.MainLedger.AddConstantProductExchange(cped, tick)
}

// NOTE LEFT HERE --- nov. 4 (need this to pass the exchange addr into a traders "make transaction"
func (s *Simulation) GetExchange(symIn, symOut string) 