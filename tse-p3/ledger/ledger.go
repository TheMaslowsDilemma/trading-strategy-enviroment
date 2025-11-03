package ledger

import (
	"fmt"
	"math/rand"
	"tse-p3/exchanges"
	"tse-p3/wallets"
)

type Ledger struct {
	Wallets		map[Addr]wallets.Wallet
	Exchanges	map[Addr]exchanges.ConstantProductExchange
}

func (l Ledger) String() string {
	var (
		wc	uint
		ec	uint
	)

	wc = 0
	ec = 0
	for _, _ = range l.Wallets {
		wc += 1
	}
	for _, _ = range l.Exchanges {
		ec += 1
	}
	return fmt.Sprintf("{ wallet-count: %v, exchange-count: %v }", wc, ec)
}

func CreateLedger() Ledger {
	return Ledger {
		Wallets: make(map[Addr]wallets.Wallet),
		Exchanges: make(map[Addr]exchanges.ConstantProductExchange),
	}
}

func RandomAddr() Addr {
	return Addr(uint64(rand.Uint32()) << 32 | uint64(rand.Uint32()))
}

func (l Ledger) AddConstantProductExchange(cd exchanges.CpeDescriptor) Addr {
	var (
		addr	Addr
		cpe	exchanges.ConstantProductExchange
	)

	addr = RandomAddr()
	cpe = exchanges.CreateConstantProductExchange(cd)
	l.Exchanges[addr] = cpe
	return addr
}

func (l Ledger) AddWallet(wd wallets.WalletDescriptor) Addr {
	var (
		addr Addr
		wlt  wallets.Wallet
	)
	
	addr = RandomAddr()
	wlt  = wallets.CreateWallet(wd)
	l.Wallets[addr] = wlt
	
	return addr
}

func (l Ledger) GetWallet(addr Addr) wallets.Wallet {
	return l.Wallets[addr]
}

func (l Ledger) GetExchange(addr Addr) exchanges.ConstantProductExchange {
	return l.Exchanges[addr]
}

// NOTE this is really pseudo merge, it should eventually support deletes
func (l Ledger) Merge(feat Ledger) {
	var (
		featwlt wallets.Wallet
		featexg exchanges.ConstantProductExchange
		addr	Addr
	)

	// merge wallet subledger - NOTE we do not do any deletes
	for addr, featwlt = range feat.Wallets {
		l.Wallets[addr] = featwlt
	}

	// merge exchange subledger
	for addr, featexg = range feat.Exchanges {
		l.Exchanges[addr] = featexg
	}
}
