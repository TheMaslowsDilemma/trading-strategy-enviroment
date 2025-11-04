package txs

import (
	"fmt"
	"tse-p3/wallets"
	"tse-p3/exchanges"
	"tse-p3/ledger"
	"github.com/holiman/uint256"
)

type CpeSwap struct {
	SymbolIn		string
	SymbolOut		string
	AmountIn		*uint256.Int
	AmountMinOut	*uint256.Int
	WalletAddr		ledger.Addr
	ExchangeAddr	ledger.Addr
	Notify			func (res TxResult)
}

// -- returns a partial ledger with values to update -- //
func (tx CpeSwap) Apply(tick uint64, l ledger.Ledger) (ledger.Ledger, error) {
	var (
		exg	exchanges.ConstantProductExchange
		wlt	wallets.Wallet
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
	fmt.Printf("wallet: %v places trade on exchange %v at price %v lmod: %v", wlt, exg, price, lmod)
	
	return lmod, nil
}
