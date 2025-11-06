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
	NeedsWallet		bool
	ExchangeAddr	ledger.Addr
	Notifier		func (res TxResult)
}

func (tx CpeSwap) Notify(res TxResult) {
	tx.Notifier(res)
}

// -- returns a partial ledger with values to update -- //
func (tx CpeSwap) Apply(tick uint64, l ledger.Ledger) (ledger.Ledger, error) {
	var (
		exg		exchanges.ConstantProductExchange
		wlt		wallets.Wallet
		lmod	ledger.Ledger
		price	*uint256.Int
		//amtout	*uint256.Int
	)

	// modified ledger
	lmod = ledger.CreateLedger()

	// --- We might not have a receiving wallet yet, so make it if needed --- ///
	if tx.NeedsWallet {
		tx.WalletAddr = lmod.AddWallet(wallets.WalletDescriptor {
			Amount: 0,
			Symbol: tx.SymbolOut,
		})
		wlt = lmod.GetWallet(tx.WalletAddr)
	} else {
		wlt = l.GetWallet(tx.WalletAddr)
	}

	// retrieve entities involved in swap
	exg = l.GetExchange(tx.ExchangeAddr).Clone()
	if exg.Auditer == nil {
		return lmod, fmt.Errorf("no exchange found %v <-> %v", tx.SymbolIn, tx.SymbolOut)
	}
	price = exg.SpotPriceA()

	// NOTE this is a test modification 
	exg.ReserveA.Amount.Sub(exg.ReserveA.Amount, tx.AmountIn)
	lmod.Exchanges[tx.ExchangeAddr] = exg
	// just print for now
	fmt.Printf("wallet: %v executed trade on exchange %v at price %v lmod: %v\n", wlt, exg, price, lmod)
	
	return lmod, nil
}
