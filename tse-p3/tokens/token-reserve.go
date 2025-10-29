package tokens

import (
	"fmt"
	"github.com/holiman/uint256"
	"github.com/cespare/xxhash"
)

type TokenReserve struct {
	Amount	*uint256.Int
	Symbol	string
}

func CreateTokenReserve(amt uint64, symb string) TokenReserve {
	return TokenReserve {
		Amount: uint256.NewInt(amt),
		Symbol: symb,
	}
}

// --- Ledger Item Implementation --- //
func (tkr TokenReserve) Clone() TokenReserve {
	return TokenReserve {
		Amount: tkr.Amount.Clone(),
		Symbol: tkr.Symbol, // NOTE we might need to clone this
	}
}

func (tkr TokenReserve) String() string {
	return fmt.Sprintf("{ amt: %v; symbol: %v }", tkr.Amount, tkr.Symbol)
}

func (tkr TokenReserve) Hash() uint64 {
	return xxhash.Sum64([]byte(tkr.String()))
}
