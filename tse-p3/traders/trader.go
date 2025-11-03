package traders

import (
	"fmt"
	"tse-p3/globals"
	"tse-p3/ledger"
	"github.com/holiman/uint256"
	"github.com/cespare/xxhash"
)

type Trader struct {
	Id		uint64
	Wallets 	map[uint64] ledger.Addr
	LastNetworth	*uint256.Int
}

func (t Trader) GetWallet(sym string) (ledger.Addr, error) {
	var (
		key	uint64
		addr	ledger.Addr
	)
	
	key = xxhash.Sum64([]byte(sym))
	addr = t.Wallet[key]

	if addr == 0 {
		return 0, fmt.Errorf("trader has no wallet for token %v", sym)
	}
	return addr, nil
}

func (t *Trader) AddWallet(sym string, addr ledger.Addr) {
	var setaddr = t.SetWallet(sym, addr, false)
	if setaddr != addr {
		fmt.Printf("trader already has wallet %v\n", sym)
	}
}

func (t *Trader) SetWallet(sym string, addr ledger.Addr, override bool) ledger.Addr {
var (
		key 	uint64
		crnt	ledger.Addr
	)

	key = xxhash.Sum64([]byte(sym))	
	crnt = t.Wallets[key]
	
	if crnt != 0 & !override {
		return crnt
	}

	t.Wallets[key] = addr
	return addr
}

func (t Trader) GetNetworth(rateProvider ledger.RateProvider, walletProvider ledger.WalletProvider) {
	var (
		networth	*uint256.Int
		waddr		ledger.Addr
		wlt		wallets.Wallet
		rate		*uint256.Int
		worth		*uint256.Int
		err		error
	)

	networth = uint256.NewInt(0)
	worth	 = uint256.NewInt(0)
	for _, waddr = range t.Wallets {
		wlt, err = walletProvider(waddr)
		if err != nil {
			fmt.Printf("trader has invalid wallet: %v\n", waddr)
			continue
		}
		rate, err = rateProvider(wlt.Reserve.Symbol, globals.BaseCurrencySymbol)
		if err != nil {
			// this token is disconnected from the base currency
			continue
		}
		networth.AddI(worth.Mul(wlt.Reserve.Amount, rate))
	}
	return networth
}
