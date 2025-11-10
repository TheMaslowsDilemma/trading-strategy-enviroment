package txs

import (
	"fmt"
	"tse-p3/wallets"
	"tse-p3/exchanges"
	"tse-p3/traders"
	"tse-p3/ledger"
	"github.com/holiman/uint256"
)

type CpeSwapDescriptor struct {
	SymbolIn	string
	SymbolOut	string
	AmountIn	*uint256.Int
	AmountMinOut	*uint256.Int
	Notifier	func (res TxResult)
}

type CpeSwap struct {
	SymbolIn		string
	SymbolOut		string
	AmountIn		*uint256.Int
	AmountMinOut		*uint256.Int
	Trader			*traders.Trader
	ExchangeAddr		ledger.Addr
	Notifier		func (res TxResult)
}

func (tx CpeSwap) Notify(res TxResult) {
	tx.Notifier(res)
}

// -- returns a partial ledger with values to update -- //
func (tx CpeSwap) Apply(tick uint64, l ledger.Ledger) (ledger.Ledger, error) {
	var (
		exg			exchanges.ConstantProductExchange
		pyr_wlt			wallets.Wallet
		rcv_wlt			wallets.Wallet
		pyr_wlt_addr		ledger.Addr
		rcv_wlt_addr		ledger.Addr
		ledger_delta		ledger.Ledger
		pyr_wlt_exists		bool
		rcv_wlt_exists		bool
		amt_out			*uint256.Int
		price			float64
	)

	ledger_delta = ledger.CreateLedger()
	(&ledger_delta).Merge(l)

	// --- First Find the Exchange --- //
	exg = ledger_delta.GetExchange(tx.ExchangeAddr).Clone()
	if exg.Auditer == nil {
		return ledger_delta, fmt.Errorf("no exchange found %v <-> %v", tx.SymbolIn, tx.SymbolOut)
	}

	// --- Second Find the Payor and Recipient Wallet Addr --- //
	pyr_wlt_addr, pyr_wlt_exists = tx.Trader.GetWalletAddr(tx.SymbolIn)
	if !pyr_wlt_exists {
		return ledger_delta, fmt.Errorf("payor wallet DNE.")
	}

	// --- If the trader has no wallet to recieve yet then make one --- //
	rcv_wlt_addr, rcv_wlt_exists = tx.Trader.GetWalletAddr(tx.SymbolOut)
	if !rcv_wlt_exists {
		rcv_wlt_addr = (&ledger_delta).AddWallet(wallets.WalletDescriptor {
			Amount: 0,
			Symbol: tx.SymbolOut,
		})
		tx.Trader.AddWallet(tx.SymbolOut, rcv_wlt_addr)
	}

	pyr_wlt = ledger_delta.GetWallet(pyr_wlt_addr).Clone()
	rcv_wlt = ledger_delta.GetWallet(rcv_wlt_addr).Clone()

	if pyr_wlt.Reserve.Amount.Lt(tx.AmountIn) {
		return ledger_delta, fmt.Errorf("insufficient funds.")
	}
	if tx.SymbolIn == exg.ReserveA.Symbol {
		amt_out = exg.SwapAForB(tx.AmountIn)
		exg.ReserveB.Amount.Sub(exg.ReserveB.Amount, amt_out)
		exg.ReserveA.Amount.Add(exg.ReserveA.Amount, tx.AmountIn)
	} else {
		amt_out = exg.SwapBForA(tx.AmountIn)
		exg.ReserveA.Amount.Sub(exg.ReserveA.Amount, amt_out)
		exg.ReserveB.Amount.Add(exg.ReserveB.Amount, tx.AmountIn)
	}

	// Update the Traders Wallets
	pyr_wlt.Reserve.Amount.Sub(pyr_wlt.Reserve.Amount, tx.AmountIn)
	rcv_wlt.Reserve.Amount.Add(rcv_wlt.Reserve.Amount, amt_out)

	price = exg.SpotPriceA()
	exg.Auditer.Audit(price, tick)
	fmt.Printf("[%v] price: %v\n", tick, price)
	// --- Finally Write changes to our delta ledger --- //
	ledger_delta.Exchanges[tx.ExchangeAddr] = exg
	ledger_delta.Wallets[rcv_wlt_addr] = rcv_wlt
	ledger_delta.Wallets[pyr_wlt_addr] = pyr_wlt
	
	return ledger_delta, nil
}
