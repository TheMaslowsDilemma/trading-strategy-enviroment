package simulation

import (
	"fmt"
	"tse-p3/users"
	"tse-p3/traders"
	"tse-p3/exchanges"
	"tse-p3/ledger"
	"tse-p3/wallets"
	"tse-p3/globals"
	"tse-p3/transactions"
	"github.com/holiman/uint256"
)

// --- Entity Functionality --- //
func (s *Simulation) PlaceUserSwap(userkey uint64, from, to string, amount uint64) {
	usr := s.Users[userkey]
	eaddr := s.ExchangeDirectory[getExchangeKey(from,to)]
	
	swaptx := txs.CpeSwap {
		SymbolIn: from,
		SymbolOut: to,
		AmountIn: uint256.NewInt(1000),
		AmountMinOut: uint256.NewInt(0),
		Trader: s.Traders[usr.TraderId],
		ExchangeAddr: eaddr,
		Notifier: Notificationator,
	}

	s.placeTx(swaptx)
}

func Notificationator(res txs.TxResult) {
	fmt.Printf("tx result: %v\n", res)
}

func (s *Simulation) AddUser(name string, pubkey uint64) {
	var (
		trdr	*traders.Trader
		usr		users.User
		wd		wallets.WalletDescriptor
		waddr	ledger.Addr
	)

	trdr = traders.CreateTrader()
	wd = wallets.WalletDescriptor {
		Amount: globals.UserStartingBalance,
		Symbol: globals.TSECurrencySymbol,
	}

	waddr = s.addWallet(wd) // Add wallet to ledger
	trdr.AddWallet(wd.Symbol, waddr) // Add wallet address to trader
	s.addTrader(trdr) // Add Trader to simulation

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

func (s *Simulation) addTrader(t *traders.Trader) {
	s.Traders[t.Id] = t
}

func (s *Simulation) addWallet(wd wallets.WalletDescriptor) ledger.Addr {
	return (&s.ScndLedger).AddWallet(wd) // NOTE: we add to back ledger because that handles all updates! CAUTION multiple writes
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