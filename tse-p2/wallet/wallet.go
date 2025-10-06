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
	if  w.ContainsReserve(raddr) {
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
        var rs []ledger.LedgerAddr = make([]ledger.LedgerAddr, len(w.Reserves))
	copy(w.Reserves, rs)
        return Wallet {
                TraderId: w.TraderId,
                Reserves: rs,
	}
}

func (w Wallet) ContainsReserve(raddr ledger.LedgerAddr) bool {
    var i int;
    for i = 0; i < len(w.Reserves); i++ {
        if w.Reserves[i] == raddr {
            return true
        }
    }
    return false
}

