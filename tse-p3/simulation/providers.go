package simulation

import (
	"fmt"
	"tse-p3/candles"
	"tse-p3/wallets"
	"tse-p3/exchanges"
	"tse-p3/ledger"
	"tse-p3/traders"
	"tse-p3/transactions"
	"github.com/holiman/uint256"
)

func (s *Simulation) placeTx(tx txs.Tx) bool {
	return s.MemoryPool.Push(tx)
}

func (s Simulation) GetWallet(waddr ledger.Addr) (wallets.Wallet, error) {
	var (
		wlt	wallets.Wallet
	)

	wlt = s.SecondaryLedger.GetWallet(waddr)
	// NOTE this seems like less than ideal way to check if the wallet exists
	if wlt.Reserve.Amount == nil {
		return wallets.Wallet{}, fmt.Errorf("no wallet exists for addr: %v", waddr)
	}

	return wlt, nil
}

func (s Simulation) GetPrice(symbol, inTermsOf string) (float64, error) {
	var (
		exkey	uint64
		exaddr	ledger.Addr
		exg		exchanges.ConstantProductExchange
	)
	exkey = getExchangeKey(symbol, inTermsOf)
	exaddr = s.ExchangeDirectory[exkey]
	if exaddr == 0 {
		return 0, fmt.Errorf("no direct exchange exists for %v <-> %v", symbol, inTermsOf)
	}
	
	exg = s.PrimaryLedger.GetExchange(exaddr)
	if exg.Auditer == nil {
		return 0, fmt.Errorf("exchange is malformed or DNE: %v", exaddr)
	}

	if symbol == exg.ReserveA.Symbol {
		return exg.SpotPriceA(), nil
	}

	return exg.SpotPriceB(), nil
}

func (s Simulation) GetCandles(symbolA, symbolB string) []candles.Candle {
		var (
		exkey	uint64
		exaddr	ledger.Addr
		exg		exchanges.ConstantProductExchange
	)
	exkey 	= getExchangeKey(symbolA, symbolB)
	exaddr 	= s.ExchangeDirectory[exkey]
	exg 	= s.PrimaryLedger.Exchanges[exaddr]

	if exg.Auditer == nil {
		return []candles.Candle{}
	}
	
	return exg.Auditer.GetCandles()
}

func (s Simulation) GetNetworth(traderKey uint64) (*uint256.Int, error) {
	var (
		tr	*traders.Trader
	)

	tr = s.Traders[traderKey]
	if tr.Id == 0 { 
		return nil, fmt.Errorf("no trader exists for key: %v", traderKey)
	}
	return tr.GetNetworth(s.GetPrice, s.GetWallet), nil
}
