package tokens

import (
	"fmt"
	"tse-p3/globals"
	"github.com/holiman/uint256"
	"github.com/cespare/xxhash"
)

type TokenReserve struct {
	Amount	*uint256.Int
	Symbol	string
}

type Descriptor struct {
	Amount	uint64
	Symbol	string
}

func CreateTokenReserve(td Descriptor) TokenReserve {
	var scaled *uint256.Int = uint256.NewInt(td.Amount)
	scaled.Mul(scaled, globals.TokenScaleFactor)
	return TokenReserve {
		Amount: scaled,
		Symbol: td.Symbol,
	}
}

func (tkr *TokenReserve) Merge(feat TokenReserve) {
	tkr.Amount = feat.Amount
	tkr.Symbol = feat.Symbol
}

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
