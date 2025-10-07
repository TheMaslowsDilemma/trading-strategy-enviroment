package token

import (
	"fmt"
        "tse-p2/ledger"
        "crypto/sha256"
)

type TokenReserve struct {
	Symbol	string
	Amount	uint64
}

func (lp TokenReserve) Copy() ledger.LedgerItem {
	return TokenReserve {
		Symbol: lp.Symbol,
		Amount: lp.Amount,
	}
}

func (lp TokenReserve) String() string {
	return fmt.Sprintf("{ sym: \"%v\", amt: %v }", lp.Symbol, lp.Amount)
}

func (lp TokenReserve) Hash() [sha256.Size]byte {
	return sha256.Sum256([]byte(lp.String()))
}

// Note: This is returning a value, NOT casting the memory region at li
func TkrFromLedgerItem(li ledger.LedgerItem) (*TokenReserve, error) {
    var (
        tkr     TokenReserve
        ok      bool
    )

    if li == nil {
        return nil, fmt.Errorf("cannot cast tkr from nil ledger item.")
    }

    if tkr, ok = (li).(TokenReserve); ok {
        return &tkr, nil
    }

    return nil, fmt.Errorf("cannot cast tkr from non-tkr ledger item.")
}
