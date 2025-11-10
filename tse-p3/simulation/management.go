package simulation

import (
	"fmt"
	"tse-p3/users"
	"tse-p3/bots"
	"tse-p3/traders"
	"tse-p3/exchanges"
	"tse-p3/ledger"
	"tse-p3/wallets"
	"tse-p3/globals"
	"tse-p3/transactions"
	"tse-p3/strategy"
	"github.com/holiman/uint256"
)

// --- Entity Functionality --- //
func (s *Simulation) PlaceUserSwap(userkey uint64, from, to string, amount uint64) {
	usr := s.Users[userkey]
	eaddr := s.ExchangeDirectory[getExchangeKey(from,to)]
	amt_in := uint256.NewInt(amount)

	swaptx := txs.CpeSwap {
		SymbolIn: from,
		SymbolOut: to,
		AmountIn: amt_in,
		AmountMinOut: uint256.NewInt(0),
		Trader: s.Traders[usr.TraderId],
		ExchangeAddr: eaddr,
		Notifier: Notificationator(usr.Name),
	}

	s.placeTx(swaptx)
}

func (s *Simulation) PlaceBotSwap(botkey uint64, dscr txs.CpeSwapDescriptor) {
	bot := s.Bots[botkey]
	eaddr := s.ExchangeDirectory[getExchangeKey(dscr.SymbolIn, dscr.SymbolOut)]
	
	swaptx := txs.CpeSwap {
			SymbolIn: dscr.SymbolIn,
			SymbolOut: dscr.SymbolOut,
			AmountIn: dscr.AmountIn,
			AmountMinOut: dscr.AmountMinOut,
			ExchangeAddr: eaddr,
			Trader: bot.Trader,
			Notifier: dscr.Notifier,
	}
	s.placeTx(swaptx)
}


func Notificationator(name string) func (txs.TxResult) {
	return func (res txs.TxResult) {
		fmt.Printf("[%v] tx result: %v\n", name, res)
	}
}

func (s *Simulation) AddUser(name string, pubkey uint64) {
	var (
		trdr	*traders.Trader
		usr	users.User
		wd	wallets.WalletDescriptor
		waddr	ledger.Addr
	)

	trdr = traders.CreateTrader()
	wd = wallets.WalletDescriptor {
		Amount: globals.UserStartingBalance,
		Symbol: globals.USDSymbol,
	}

	waddr = s.AddWallet(wd) // Add wallet to ledger
	trdr.AddWallet(wd.Symbol, waddr) // Add wallet address to trader
	s.AddTrader(trdr) // Add Trader to simulation


	usr = users.User {
		Name: name,
		PublicKey: pubkey,
		TraderId: trdr.Id,
	}
	s.Users[pubkey] = usr 
}

func (s *Simulation) AddBot(name string, strat strategies.Strategy) uint64 {
		var (
		trdr	*traders.Trader
		bot	bots.Bot
		wd	wallets.WalletDescriptor
		waddr	ledger.Addr
	)

	trdr = traders.CreateTrader()
	wd = wallets.WalletDescriptor {
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
	return (&s.SecondaryLedger).AddWallet(wd) // NOTE: we add to back ledger because that handles all updates! CAUTION multiple writes
}

func (s *Simulation) AddExchange(cd exchanges.CpeDescriptor, tick uint64) {
	var eaddr ledger.Addr
	var dirKeyForward, dirKeyBackward uint64
	// NOTE consider just sorting the symbols in the "getExchangeKey" func
	// so both forward and backward return the same key
	dirKeyForward = getExchangeKey(cd.SymbolA, cd.SymbolB)
	dirKeyBackward = getExchangeKey(cd.SymbolB, cd.SymbolA)
	eaddr = s.PrimaryLedger.AddConstantProductExchange(cd, tick)
	s.ExchangeDirectory[dirKeyForward] = eaddr
	s.ExchangeDirectory[dirKeyBackward] = eaddr
}

func (s *Simulation) GetExchange(symIn, symOut string) exchanges.ConstantProductExchange{
	var eaddr uint64 = getExchangeKey(symIn, symOut)
	return s.PrimaryLedger.GetExchange(ledger.Addr(eaddr))
}
