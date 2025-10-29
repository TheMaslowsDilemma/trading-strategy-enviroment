package wallets

import (
	"fmt"
	"tse-p3/tokens"
	"github.com/cespare/xxhash"
)

type Wallet struct {
	Reserve	tokens.TokenReserve
}

func CreateWallet(amt uint64, symb string) Wallet {
	return Wallet {
		Reserve: tokens.CreateTokenReserve(amt, symb),
	}
}

// Ledger Item Implementation //
// TODO Merge() operation

func (wlt Wallet) Clone() {
	return Wallet {
		Reserve: wlt.Reserve.Clone(),
	}
}

func (wlt Wallet) String() {
	return fmt.Sprintf("{ reserve: %v }", wlt.Reserve)
}

func (wlt Wallet) Hash() uint64 {
	return xxhash.Sum64([]byte(wlt.String()))
}
