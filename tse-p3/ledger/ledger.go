package ledger

import (
	"fmt"
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

func NewLedgerWithEmitter() *Ledger {
	return &Ledger{
		Wallets:		make(map[Addr]wallets.Wallet),
		Exchanges:		make(map[Addr]exchanges.ConstantProductExchange),
		EmitManager:	NewEmitterManager(),
	}
}

// -------------------helpers--------------------------

func (l *Ledger) AddWallet(name string, wd wallets.WalletDescriptor) Addr {
	addr := globals.Rand64()
	wallet := wallets.CreateWallet(wd)
	l.Wallets[addr] = wallet

	if l.EmitManager != nil {
		l.EmitManager.AddSource(name, addr, EntityWallet)
	}
	return addr
}

func (l *Ledger) AddConstantProductExchange(name string, cd exchanges.CpeDescriptor, tick uint64) Addr {
	addr := globals.Rand64()
	cpe := exchanges.CreateConstantProductExchange(cd, tick)
	l.Exchanges[addr] = cpe

	if l.EmitManager != nil {
		l.EmitManager.AddSource(name, addr, EntityExchange)
	}
	return addr
}

func (l *Ledger) SearchSources(name string) []SearchResult {
	if l.EmitManager != nil {
		return l.EmitManager.SearchSources(name)
	}
	return []SearchResult{}
}

func (l *Ledger) GetWallet(addr Addr) (wallets.Wallet, bool) {
	w, ok := l.Wallets[addr]
	return w, ok
}

func (l *Ledger) GetExchange(addr Addr) (exchanges.ConstantProductExchange, bool) {
	e, ok := l.Exchanges[addr]
	return e, ok
}

func (l *Ledger) Merge(delta Ledger) {
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

func (l *Ledger) MergeAndEmit(delta Ledger) {
	if l.EmitManager == nil {
		l.Merge(delta) // fallback to non emitter merge
		return
	}

	// --- Wallets ---
	for addr, newWallet := range delta.Wallets {
		oldWallet, exists := l.Wallets[addr]

		if !exists || oldWallet.Hash() != newWallet.Hash() {
			// Update state
			l.Wallets[addr] = newWallet

			// Emit change
			dsrc := l.EmitManager.AddSource(addr, EntityWallet) // ensures source exists
			dsrc.Emit(map[string]any{
				"name": dsrc.Name,
				"address": addr,
				"balance": newWallet.Reserve.Amount.String(),
				"token":   newWallet.Reserve.Token.Symbol,
			}, EntityWallet)
		}
	}

	// --- Exchanges ---
	for addr, newEx := range delta.Exchanges {
		oldEx, exists := l.Exchanges[addr]

		if !exists || oldEx.Hash() != newEx.Hash() {
			// Update state
			l.Exchanges[addr] = newEx

			// Emit spot price of token A (or both if you prefer)
			priceA := newEx.SpotPriceA()
			priceB := newEx.SpotPriceB()

			dsrc := l.EmitManager.AddSource(addr, EntityExchange)
			dsrc.Emit(map[string]any{
				"name":		dsrc.Name,
				"address":  addr,
				"priceA":   priceA.String(),
				"priceB":   priceB.String(),
				"reserveA": newEx.ReserveA.Amount.String(),
				"reserveB": newEx.ReserveB.Amount.String(),
				"tokenA":   newEx.ReserveA.Symbol,
				"tokenB":   newEx.ReserveB.Symbol,
			}, EntityExchange)
		}
	}
}