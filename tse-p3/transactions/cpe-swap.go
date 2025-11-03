package txs

import (
	"fmt"
	"github.com/holiman/uint256"
	"tse-p3/ledger"
)

type CpeSwap struct {
	SymbolIn	string
	SymbolOut	string
	AmountIn	*uint256.Int
	AmountMinOut	*uint256.Int
	WalletAddr	ledger.Addr
	ExchangeAddr	ledger.Addr
	Notify		func (res ledger.TxResult)
}

// -- returns a partial ledger with values to update -- //
func (tx CpeSwap) Apply(tick uint64, l ledger.Ledger) (ledger.Ledger, error) {
	var (
		exg	ConstantProductExchange
		wlt	wallet.Wallet
		lmod	ledger.Ledger
		price	*uint256.Int
//		amtout  *uint256.Int
	)
	
	// modified ledger
	lmod = ledger.CreateLedger()

	// retrieve entities involved in swap
	wlt = l.GetWallet(tx.WalletAddr)
	exg = l.GetExchange(tx.ExchangeAddr)
	price = exg.SpotPriceA()

	// just print for now
	fmt.Printf("wallet: %v places trade on exchange %v at price %v", wlt, exg, price)
	
	return l, nil
}
