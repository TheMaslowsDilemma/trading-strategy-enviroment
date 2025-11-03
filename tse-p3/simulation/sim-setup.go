package simulation

import (
	"tse-p3/exchanges"
	"tse-p3/ledger"
	"tse-p3/wallets"
)

func (s *Simulation) AddUser() {
	// Add User
	// Add Trader
	// Add Wallet
}

func (s *Simulation) AddWallet(wd wallets.WalletDescriptor) ledger.Addr {
	return s.MainLedger.AddWallet(wd)
}

func (s *Simulation) AddExchange(cped exchanges.CpeDescriptor, tick uint64) ledger.Addr {
	return s.MainLedger.AddConstantProductExchange(cped, tick)
}