package ledger

import (
	"fmt"
	"strconv"
	"tse-p3/exchanges"
	"tse-p3/wallets"
	"tse-p3/globals"
)


type Ledger struct {
	Wallets		map[Addr]wallets.Wallet
	Exchanges	map[Addr]exchanges.ConstantProductExchange
	EmitManager	*EmitterManager
}

func (l *Ledger) String() string {
	return fmt.Sprintf("{wallets: %d, exchanges: %d}", len(l.Wallets), len(l.Exchanges))
}

// Constructors
func NewLedger() *Ledger {
	return &Ledger{
		Wallets:	make(map[Addr]wallets.Wallet),
		Exchanges:	make(map[Addr]exchanges.ConstantProductExchange),
	}
}

func Clone(to_clone Ledger) *Ledger {
	var clone *Ledger = NewLedger()
	for addr, w := range to_clone.Wallets {
		clone.Wallets[addr] = w
	}
	for addr, e := range to_clone.Exchanges {
			clone.Exchanges[addr] = e
	}
	return clone
}

func NewLedgerWithEmitter() *Ledger {
	return &Ledger{
		Wallets:		make(map[Addr]wallets.Wallet),
		Exchanges:		make(map[Addr]exchanges.ConstantProductExchange),
		EmitManager:	NewEmitterManager(),
	}
}

// -------------------helpers--------------------------

func (l *Ledger) AddWallet(wd wallets.WalletDescriptor) Addr {
	addr := Addr(globals.Rand64())
	wallet := wallets.CreateWallet(wd)
	l.Wallets[addr] = wallet

	if l.EmitManager != nil {
		l.EmitManager.AddSource(wallet.Name, addr, Wallet_t)
	}
	return addr
}

func (l *Ledger) AddConstantProductExchange(cd exchanges.CpeDescriptor, tick uint64) Addr {
	addr := Addr(globals.Rand64())
	cpe := exchanges.CreateConstantProductExchange(cd, tick)
	name := fmt.Sprintf("%v <=> %v",
		cpe.ReserveA.Symbol,
		cpe.ReserveB.Symbol,
	)
	l.Exchanges[addr] = cpe

	if l.EmitManager != nil {
		l.EmitManager.AddSource(name, addr, Exchange_t)
	}
	return addr
}

func (l *Ledger) SearchSources(name string) []SearchResult {
	if l.EmitManager != nil {
		return l.EmitManager.SearchSources(name)
	}
	return []SearchResult{}
}

func (l *Ledger) GetWallet(addr Addr) wallets.Wallet {
	return l.Wallets[addr]
}

func (l *Ledger) GetExchange(addr Addr) exchanges.ConstantProductExchange {
	return l.Exchanges[addr]
}

func (l *Ledger) Merge(delta *Ledger) {
	for addr, w := range delta.Wallets {
		if current, exists := l.Wallets[addr]; !exists || current.Hash() != w.Hash() {
			l.Wallets[addr] = w
		}
	}
	for addr, e := range delta.Exchanges {
		if current, exists := l.Exchanges[addr]; !exists || current.Hash() != e.Hash() {
			l.Exchanges[addr] = e
		}
	}
}

func (l *Ledger) MergeAndEmit(tick uint64, delta *Ledger) {
	if l.EmitManager == nil {
		l.Merge(delta) // fallback to non emitter merge
		return
	}

	// --- Wallets ---
	for addr, newWallet := range delta.Wallets {
		// Update state
		l.Wallets[addr] = newWallet

		// Emit change
		dsrc := l.EmitManager.AddSource(newWallet.Name, addr, Wallet_t) // ensures source exists
		dsrc.Emit(tick,
			map[string]any {
				"name": dsrc.Name,
				"address": strconv.FormatUint(uint64(addr), 10),
				"balance":	newWallet.Reserve.Amount.String(),
				"symbol":	newWallet.Reserve.Symbol,
			},Wallet_t)
		
	}

	// --- Exchanges ---
	for addr, newEx := range delta.Exchanges {

		// Update state
		l.Exchanges[addr] = newEx

		// Emit candles 
		dsrc := l.EmitManager.AddSource(
			fmt.Sprintf(
				"%v <=> %v",
				newEx.ReserveA.Symbol,
				newEx.ReserveB.Symbol,
			),
			addr,
			Exchange_t,
		)
		dsrc.Emit(tick, map[string]any{
			"address":	 strconv.FormatUint(uint64(addr), 10),
			"candles":	newEx.Auditer.GetCandles(),
			"reserveA":	newEx.ReserveA.Amount.String(),
			"reserveB":	newEx.ReserveB.Amount.String(),
			"tokenA":	newEx.ReserveA.Symbol,
			"tokenB":	newEx.ReserveB.Symbol,
		}, Exchange_t)

	}
}