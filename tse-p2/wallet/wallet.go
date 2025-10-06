package wallet

import (
	"fmt"
	"tse-p2/ledger"
	"crypto/sha256"
)

type Wallet struct { 
	TokenA	uint64
	TokenB	uint64
}

func (w Wallet) String() string {
	return fmt.Sprintf("{ TokenA: %v, TokenB: %v }", w.TokenA, w.TokenB)
}

func (w Wallet) Hash() [sha256.Size]byte {
	return sha256.Sum256([]byte(w.String()))
}

func (w Wallet) Copy() ledger.LedgerItem {
	return Wallet {
		TokenA: w.TokenA,
		TokenB: w.TokenB,
	}
}