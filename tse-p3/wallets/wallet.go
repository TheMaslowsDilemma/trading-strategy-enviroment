package wallets

import (
	"fmt"
	"tse-p3/tokens"
	"github.com/cespare/xxhash"
)

type Wallet struct {
	Reserve	tokens.TokenReserve
}

type WalletDescriptor struct {
	Amount uint64
	Symbol string
}

func CreateWallet(wd WalletDescriptor) Wallet {
	return Wallet {
		Reserve: tokens.CreateTokenReserve(wd.Amount, wd.Symbol),
	}
}

// Ledger Item Implementation //
// TODO Merge() operation

func (wlt Wallet) Clone() Wallet  {
	return Wallet {
		Reserve: wlt.Reserve.Clone(),
	}
}

func (wlt Wallet) String() string {
	return fmt.Sprintf("{ reserve: %v }", wlt.Reserve)
}

func (wlt Wallet) Hash() uint64 {
	return xxhash.Sum64([]byte(wlt.String()))
}
