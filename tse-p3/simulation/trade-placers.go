package simulation

import (
	"fmt"
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

// --- User Swap  --- //
func (s *Simulation) PlaceUserSwap(trader_id uint64, from, to string, amount uint64) {
	eaddr := s.ExchangeDirectory[globals.GetExchangeKey(from,to)]
	amt_in := uint256.NewInt(amount)

	swaptx := txs.CpeSwap {
		SymbolIn: from,
		SymbolOut: to,
		AmountIn: amt_in,
		AmountMinOut: uint256.NewInt(1),
		Trader: s.Traders[trader_id],
		ExchangeAddr: eaddr,
		Notifier: Notificationator(usr.Name),
	}

	s.placeTx(swaptx)
}


// --- Bot Swap  --- //
func (s *Simulation) PlaceBotSwap(botkey uint64, dscr txs.CpeSwapDescriptor) {
	bot := s.Bots[botkey]
	eaddr := s.ExchangeDirectory[globals.GetExchangeKey(dscr.SymbolIn, dscr.SymbolOut)]
	
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
