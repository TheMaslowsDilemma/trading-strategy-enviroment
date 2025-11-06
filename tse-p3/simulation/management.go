package simulation

import (
	"fmt"
	"tse-p3/users"
	"tse-p3/traders"
	"tse-p3/exchanges"
	"tse-p3/ledger"
	"tse-p3/wallets"
	"tse-p3/globals"
	//"github.com/cespare/xxhash"
)

// --- Entity Functionality --- //
// TODO add maximum slippage
func (s *Simulation) PlaceUserSwap(userkey uint64, from, to string, amount uint64) {
	usr := s.Users[userkey]
	trdr := s.Traders[usr.TraderId]
	eaddr := s.ExchangeDirectory[getExchangeKey(from,to)]
	swaptx, err := trdr.CreateSwapTx(from, to, eaddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.placeTx(swaptx)
}

// --- Adding Entities --- //

// TODO better key system 
func (s *Simulation) AddUser(name string, pubkey uint64) {
	var (
		trdr	traders.Trader
		usr		users.User
		wd		wallets.WalletDescriptor
		waddr	ledger.Addr
	)

	trdr = traders.CreateTrader()
	wd = wallets.WalletDescriptor {
		Amount: globals.UserStartingBalance,
		Symbol: globals.TSECurrencySymbol,
	}

	waddr = s.addWallet(wd)
	trdr.AddWallet(wd.Symbol, waddr)
	s.addTrader(trdr)

	usr = users.User {
		Name: name,
		PublicKey: pubkey,
		TraderId: trdr.Id,
	}
	s.Users[pubkey] = usr 
}

func (s *Simulation) AddBot() {
	// TODO
}

func (s *Simulation) addTrader(t traders.Trader) {
	s.Traders[t.Id] = t
}

func (s *Simulation) addWallet(wd wallets.WalletDescriptor) ledger.Addr {
	return s.MainLedger.AddWallet(wd)
}

func (s *Simulation) addExchange(cd exchanges.CpeDescriptor, tick uint64) {
	var eaddr ledger.Addr
	var dirKeyForward, dirKeyBackward uint64
	// NOTE consider just sorting the symbols in the "getExchangeKey" func 
	// so both forward and backward return the same key
	dirKeyForward = getExchangeKey(cd.SymbolA, cd.SymbolB)
	dirKeyBackward = getExchangeKey(cd.SymbolB, cd.SymbolA)
	eaddr = s.MainLedger.AddConstantProductExchange(cd, tick)
	s.ExchangeDirectory[dirKeyForward] = eaddr
	s.ExchangeDirectory[dirKeyBackward] = eaddr
}

func (s *Simulation) GetExchange(symIn, symOut string) exchanges.ConstantProductExchange{
	var eaddr uint64 = getExchangeKey(symIn, symOut)
	return s.MainLedger.GetExchange(ledger.Addr(eaddr))
}