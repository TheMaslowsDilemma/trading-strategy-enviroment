package simulation

import (
	"fmt"
	"tse-p3/bots"
	"tse-p3/ledger"
	"tse-p3/traders"
	"tse-p3/wallets"
	"tse-p3/globals"
	"tse-p3/strategy"
	"tse-p3/exchanges"
)

func (s *Simulation) AddBot(name string, strat strategies.Strategy) uint64 {
		var (
		trdr	*traders.Trader
		bot		bots.Bot
		wd		wallets.WalletDescriptor
		waddr	ledger.Addr
	)

	trdr = traders.CreateTrader(name)
	wd = wallets.WalletDescriptor {
		Name: fmt.Sprintf("%v:w:%v", name, globals.USDSymbol),
		Amount: globals.UserStartingBalance,
		Symbol: globals.USDSymbol,
	}

	waddr = s.AddWallet(wd)
	trdr.AddWallet(wd.Symbol, waddr)
	s.AddTrader(trdr)

	bot = bots.Bot {
		Id: globals.Rand64(),
		Name: name,
		Strategy: strat,
		Trader: trdr,
	}

	s.Bots[bot.Id] = &bot
	return bot.Id
}

func (s *Simulation) AddTrader(t *traders.Trader) {
	s.Traders[t.Id] = t
}

func (s *Simulation) AddWallet(wd wallets.WalletDescriptor) ledger.Addr {
	var waddr ledger.Addr
	s.SecondaryLock.Lock()
	waddr = s.SecondaryLedger.AddWallet(wd) // NOTE: we add to back ledger because! CAUTION multiple writes
	s.SecondaryLock.Unlock()
	return waddr
}

func (s *Simulation) AddExchange(cd exchanges.CpeDescriptor, tick uint64) {
	var (
		eaddr	ledger.Addr
		exgkey	uint64
	)

	s.PrimaryLock.Lock()
	eaddr = s.PrimaryLedger.AddConstantProductExchange(cd, tick) // Question: should this be added to the secondary?
	s.PrimaryLock.Unlock()

	// --forward direction--
	exgkey = globals.GetExchangeKey(cd.SymbolA, cd.SymbolB)
	s.ExchangeDirectory[exgkey] = eaddr

	// --backward direction--
	exgkey = globals.GetExchangeKey(cd.SymbolB, cd.SymbolA)
	s.ExchangeDirectory[exgkey] = eaddr
}