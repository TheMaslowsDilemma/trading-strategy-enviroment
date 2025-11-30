package bots

import (
	"fmt"
	"time"
	"math/rand"
	"tse-p3/globals"
	"tse-p3/candles"
	"tse-p3/ledger"
	"tse-p3/transactions"
	"tse-p3/strategy"
	"tse-p3/wallets"
	"github.com/holiman/uint256"
)

func (bot *Bot) Run(isCanceled *bool, candleProvider func(string, string) ([]candles.Candle, string), placeSwap func (botkey uint64, dscr txs.CpeSwapDescriptor), walletProvider func (ledger.Addr) (wallets.Wallet, error)) {
	var (
		tick		uint64
		cs 			[]candles.Candle
		decision	strategies.Action
		confidence	float64
		cnf_scaled	*uint256.Int
		symbol_out	string
		symbol_in	string
		amount_in	*uint256.Int
		amount_out	*uint256.Int
		waddr		ledger.Addr
		wlt			wallets.Wallet
		exists		bool
		dscr		txs.CpeSwapDescriptor
		err			error
	)

	fmt.Printf("starting [%v] run\n", bot.Name)
	for {

		var sym_a  string
		var sym_b string
		var exg_primary string

		if *isCanceled {
			return
		}

		time.Sleep(globals.BotTaskDelay)
		tick += 1

		if (bot.PendingTx) {
			continue
		}

		sym_a, sym_b = get_exchange_symbols()
		cs, exg_primary = candleProvider(sym_a, sym_b)

		if sym_a != exg_primary {
			tmp := sym_b
			sym_b = sym_a
			sym_a = tmp
		}

		decision, confidence = bot.Strategy.Decide(cs)

		if decision == strategies.Hold {
			continue
		}

		if decision == strategies.Sell {
			symbol_in	= sym_b
			symbol_out	= sym_a
		} else if decision == strategies.Buy {
			symbol_in	= sym_a
			symbol_out	= sym_b
		}

		waddr, exists = bot.Trader.GetWalletAddr(symbol_in)
		
		if !exists {
			continue
		}

		wlt, err = walletProvider(waddr)
		if err != nil {
			fmt.Printf("bot failed to get wallet: %v\n", err)
			continue // EXIT EARLY -- CORRUPT BOT
		}

		cnf_scaled = uint256.NewInt(uint64(confidence * globals.TokenScaleFactorf64))
		amount_in = cnf_scaled.Div(cnf_scaled.Mul(cnf_scaled, wlt.Reserve.Amount), globals.TokenScaleFactor)
		amount_out = uint256.NewInt(1)

		dscr = txs.CpeSwapDescriptor {
			AmountIn:	amount_in,
			SymbolIn: 	symbol_in,
			AmountMinOut:	amount_out,
			SymbolOut:	symbol_out,
			Notifier:	bot.NotificationHandler,
		}
	 	bot.PendingTx = true

		placeSwap(bot.Id, dscr)
	}
}

// TODO: create strategy to choose the exchange to trade on.
func get_exchange_symbols() (string, string) {
	var exchange_symbols []string = []string {
		globals.USDSymbol,
		globals.TSESymbol,
		"alpha",
	}

	var token_a = rand.Uint32() % 3
	var token_b = (rand.Uint32()  % 2 + token_a + 1) % 3

	return exchange_symbols[token_a], exchange_symbols[token_b]
}

