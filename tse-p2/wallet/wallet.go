package wallet

import (
	"fmt"
	"tse-p2/ledger"
	"crypto/sha256"
)

type Wallet struct {
	TraderId	uint64
	Reserves	[]ledger.LedgerAddr
}

func (w *Wallet) AddReserve(raddr ledger.LedgerAddr) error {
	if Contains(w.Reserves, raddr) {
            return fmt.Errorf("wallet already contains reserve %v", raddr)
        }
        w.Reserves = append(w.Reserves, raddr)
	return nil
}

func (w Wallet) String() string {
	return fmt.Sprintf("trader-id: %v, reserves: %v", w.TraderId, w.Reserves)
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
