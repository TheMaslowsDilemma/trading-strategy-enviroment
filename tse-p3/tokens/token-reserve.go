package tokens

import (
	"fmt"
	"tse-p3/ledger"
	"github.com/holiman/uint256"
	"github.com/cespare/xxhash"
)

type TokenReserve struct {
	Addr 	ledger.Addr
	Amount	*uint256.Int
	Symbol	string
}

func CreateTokenReserve(sym string, amt uint64) TokenReserve {
	return TokenReserve {
		Addr: 0,
		Amount: uint256.NewInt(amt),
		Symbol: sym,
	}
}

// --- Ledger Item Implementation --- //
func (tkr TokenReserve) Clone() TokenReserve {
	return TokenReserve {
		Addr: 	tkr.Addr,
		Amount: tkr.Amount.Clone(),
		Symbol: tkr.Symbol, // NOTE we might need to clone this
	}
}

func (tkr TokenReserve) String() string {
	return fmt.Sprintf("{ addr: %v; amt: %v; symbol: %v }", tkr.Addr, tkr.Amount, tkr.Symbol)
}

func (tkr TokenReserve) Hash() uint64 {
	return xxhash.Sum64([]byte(tkr.String()))
}
