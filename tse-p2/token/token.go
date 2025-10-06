package liquidity

import (
	"fmt"
	"crypto/sha256"
)

type TokenReserve struct {
	Symbol	string
	Amount	uint64
}

func (lp TokenReserve) Copy() LedgerItem {
	return TokenReserve {
		Symbol: lp.Symbol,
		Amount: lp.Amount,
	}
}

func (lp TokenReserve) String() string {
	return fmt.Sprintf("{ sym: \"%v\", amt: %v }", lp.TokenSymbol, lp.TokenCount)
}

func (lp TokenReserve) Hash() [sha256.Size]byte {
	sha256.Sum256([]byte(lp.String()))
}