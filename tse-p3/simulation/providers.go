package simulation

import (
	"fmt"
	"tse-p3/wallets"
	"tse-p3/exchanges"
	"github.com/cespare/xxhash"
	"github.com/holiman/uint256"
)

func exchangeKey(s1, s2 string) string {
	return xxhash.Sum64([]byte(fmt.Sprintf("%v:%v", s1, s2)))
}

func (s Simulation) walletProvider(waddr ledger.Addr) (wallet.Wallet, error) {
	var (
		wlt	wallet.Wallet
	)

	wlt = s.MainLedger.GetWallet(waddr)
	// NOTE this seems like less than ideal way to check if the wallet exists
	if wlt.Reserve.Amount == nil {
		return Wallet{}, fmt.Errorf("no wallet exists for addr: %v", waddr)
	}

	return wlt, nil
}

func (s Simulation) rateProvider(symbol, inTermsOf string) (*uint256.Int, error) {
	// step one
	// get the exchange key
	var (
		exkey	uint64
		exaddr	ledger.Addr
		exg	exchanges.ConstantProductExchange
	)
	exkey = exchangeKey(symbol, inTermsOf)
	exaddr = s.ExchangeDirectory[exkey]
	if exaddr == 0 {
		return nil, fmt.Errorf("no direct exchange exists for %v <-> %v", symbol, inTermsOf)
	}
	
	exg = s.Ledger.GetExchange(exaddr)
	if exg.Auditer == nil {
		return nil, fmt.Errorf("exchange is malformed or DNE: %v", exgaddr)
	}

	if symbol == exg.ReserveA.Symbol {
		return exg.SpotPriceA()
	}
	return exg.SpotPriceB()
}
