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
