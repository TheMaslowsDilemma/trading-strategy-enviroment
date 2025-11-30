package simulation

import (
	"fmt"
	"tse-p3/candles"
	"tse-p3/wallets"
	"tse-p3/exchanges"
	"tse-p3/ledger"
	"tse-p3/traders"
	"tse-p3/transactions"
	"tse-p3/globals"
)

func (s *Simulation) placeTx(tx txs.Tx) bool {
	return s.MemoryPool.Push(tx)
}

// --- Data Subscriber Logic --- //
func (s *Simulation) AddDataSubscriber(name string, addr ledger.Addr, etype ledger.EntityType, userID uint64, emitter ledger.Emit) {
	s.PrimaryLedger.EmitManager.AddSubscriber(name, addr, etype, userID, emitter)
}

func (s *Simulation) RemoveDataSubscriber(addr ledger.Addr, etype ledger.EntityType, userID uint64) {
	s.PrimaryLedger.EmitManager.RemoveSubscriber(addr, etype, userID)
}

func (s *Simulation) SearchDataSources(name string) []ledger.SearchResult {
	return s.PrimaryLedger.EmitManager.SearchSources(name)
}
// -----------------------------//



func (s *Simulation) GetWallet(waddr ledger.Addr) (wallets.Wallet, error) {
	var (
		wlt	wallets.Wallet
	)

	s.PrimaryLock.Lock()
	wlt = s.PrimaryLedger.GetWallet(waddr)
	s.PrimaryLock.Unlock()

	// NOTE this seems like less than ideal way to check if the wallet exists
	if wlt.Reserve.Amount == nil {
		return wallets.Wallet{}, fmt.Errorf("no wallet exists for addr: %v", waddr)
	}

	return wlt, nil
}

func (s *Simulation) GetPrice(symbol, inTermsOf string) (float64, error) {
	var (
		exkey	uint64
		exaddr	ledger.Addr
		exg		exchanges.ConstantProductExchange
	)

	if symbol == inTermsOf {
		return 1.0, nil
	}
	
	exkey = globals.GetExchangeKey(symbol, inTermsOf)
	exaddr = s.ExchangeDirectory[exkey]
	if exaddr == 0 {
		return 0, fmt.Errorf("no direct exchange exists for %v <-> %v", symbol, inTermsOf)
	}
	
	s.PrimaryLock.Lock()
	exg = s.PrimaryLedger.GetExchange(exaddr)
	s.PrimaryLock.Unlock()

	if exg.Auditer == nil {
		return 0, fmt.Errorf("exchange is malformed or DNE: %v", exaddr)
	}

	if symbol == exg.ReserveA.Symbol {
		return exg.SpotPriceA(), nil
	}

	return exg.SpotPriceB(), nil
}

func (s *Simulation) GetCandles(symbolA, symbolB string) ([]candles.Candle, string) {
		var (
		exkey	uint64
		exaddr	ledger.Addr
		exg		exchanges.ConstantProductExchange
	)
	exkey 	= globals.GetExchangeKey(symbolA, symbolB)
	s.PrimaryLock.Lock()
	exaddr 	= s.ExchangeDirectory[exkey]
	exg 	= s.PrimaryLedger.Exchanges[exaddr]
	s.PrimaryLock.Unlock()


	if exg.Auditer == nil {
		return []candles.Candle{}, ""
	}
	
	return exg.Auditer.GetCandles(), exg.ReserveA.Symbol
}

func (s *Simulation) GetNetworth(traderKey uint64) (float64, error) {
	var (
		tr	*traders.Trader
	)

	tr = s.Traders[traderKey]
	if tr.Id == 0 { 
		return 0.0, fmt.Errorf("no trader exists for key: %v", traderKey)
	}
	return tr.GetNetworth(s.GetPrice, s.GetWallet), nil
}

func (s *Simulation) GetExchange(symIn, symOut string) exchanges.ConstantProductExchange{
	var (
		exg_key	uint64
		exg_adr	ledger.Addr
		cpe		exchanges.ConstantProductExchange
	)
	exg_key = globals.GetExchangeKey(symIn, symOut)
	exg_adr = s.ExchangeDirectory[exg_key]

	s.PrimaryLock.Lock()
	cpe = s.PrimaryLedger.GetExchange(exg_adr)
	s.PrimaryLock.Unlock()

	return cpe
}