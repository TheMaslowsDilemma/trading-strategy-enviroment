package main

import (
	"fmt"
	"tse-p3/wallets"
	"tse-p3/exchanges"
	"tse-p3/ledger"
)

func main() {
	fmt.Println("--- Trading Stategy Environment: Part Three ---")
	
	var (
		PrimarySymbol 	string
		PrimaryAmount	uint64
		SecondarySymbol string
		SecondaryAmount uint64
		MainLedger	ledger.Ledger
	)

	PrimarySymbol = "usd"
	PrimaryAmount = 10240835
	
	SecondarySymbol = "eth"
	SecondaryAmount = 19841776

	MainLedger = ledger.CreateLedger()

	// create Exchange Descriptor
	exDscr := exchanges.CpeDescriptor {
		AmountA: PrimaryAmount,
		AmountB: SecondaryAmount,
		SymbolA: PrimarySymbol,
		SymbolB: SecondarySymbol,
	}

	exAddr := MainLedger.AddConstantProductExchange(exDscr)
	fmt.Printf("Exchange Address: %v = %v\n", exAddr, MainLedger.GetExchange(exAddr))

	wltDscr := wallets.WalletDescriptor {
		Amount: 10000,
		Symbol: PrimarySymbol,
	}
	wltAddr := MainLedger.AddWallet(wltDscr)
	fmt.Printf("Wallet Address: %v = %v\n", wltAddr, MainLedger.GetWallet(wltAddr))
}
