package wallet

import (
	"fmt"
        "tse-p2/token"
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
    var i int
    for i = 0; i < len(w.Reserves); i++ {
        if w.Reserves[i] == raddr {
            return true
        }
    }
    return false
}

func (w Wallet) GetReserveAddr(sym string, l ledger.Ledger) (ledger.LedgerAddr, error) {
    var (
        i       int
        tkr     *token.TokenReserve
        err     error
    )

    for i = 0; i < len(w.Reserves); i++ {
        tkr, err = token.TkrFromLedgerItem(l[w.Reserves[i]])
        if err != nil {
            continue; // we can ignore these errors
        }
        if tkr.Symbol == sym {
            return w.Reserves[i], nil
        }
    }
    return 0, fmt.Errorf("wallet reserve for \"%v\" DNE.", sym)
}

func WltFromLedgerItem(li ledger.LedgerItem) (*Wallet, error) {
     var (
         wlt    Wallet
         ok     bool
     )

     if li == nil {
         return nil, fmt.Errorf("cannot cast wallet from nil ledger item.")
     }

     if wlt, ok = li.(Wallet); ok {
         return &wlt, nil
     }

     return nil, fmt.Errorf("cannot cast wallet from non-wallet ledger item.")
}
