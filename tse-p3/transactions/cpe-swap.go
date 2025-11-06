package txs

import (
	"fmt"
	"tse-p3/wallets"
	"tse-p3/exchanges"
	"tse-p3/traders"
	"tse-p3/ledger"
	"github.com/holiman/uint256"
)

type CpeSwap struct {
	SymbolIn		string
	SymbolOut		string
	AmountIn		*uint256.Int
	AmountMinOut	*uint256.Int
	Trader			*traders.Trader
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
		waddr	ledger.Addr
		haswlt	bool
		wlt		wallets.Wallet
		price	*uint256.Int
		lprime  ledger.Ledger
		//amtout	*uint256.Int
	)

	lprime = ledger.CreateLedger()
	(&lprime).Merge(l)

	waddr, haswlt = tx.Trader.GetWalletAddr(tx.SymbolIn)
	if !haswlt {
		waddr = (&lprime).AddWallet(wallets.WalletDescriptor {
			Amount: 0,
			Symbol: tx.SymbolIn,
		})
		tx.Trader.AddWallet(tx.SymbolIn, waddr) // give trader the new wallet
	}
	wlt = lprime.GetWallet(waddr).Clone()
	exg = lprime.GetExchange(tx.ExchangeAddr).Clone()

	if exg.Auditer == nil {
		return lprime, fmt.Errorf("no exchange found %v <-> %v", tx.SymbolIn, tx.SymbolOut)
	}

	price = exg.SpotPriceA()

	// NOTE this is a test modification 
	exg.ReserveA.Amount.Sub(exg.ReserveA.Amount, tx.AmountIn)
	exg.Auditer.Audit(price, tick)
	lprime.Exchanges[tx.ExchangeAddr] = exg
	lprime.Wallets[waddr] = wlt
	
	return lprime, nil
}
