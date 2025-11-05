package exchanges

import (
	"fmt"
	"tse-p3/tokens"
	"tse-p3/candles"
	"tse-p3/globals"
	"github.com/holiman/uint256"
	"github.com/cespare/xxhash"
)

// NOTE the auditer is a pointer to the exchanes auditer
// considering we have a "miner" ledger and a "main" ledger
// there is a sort of duplication here. We actually want that
// duplication and must be careful not to share "auditers"
// between exchanges on the main ledger and miner ledger.
// An implication of this is that an Auditer must be merge-able
// maybe we merge based on timestampts of the underlying candles.
type ConstantProductExchange struct {
	ReserveA	tokens.TokenReserve
	ReserveB	tokens.TokenReserve
	Auditer		*candles.Auditer // consider note
}

type CpeDescriptor struct {
	AmountA	uint64
	AmountB uint64
	SymbolA string
	SymbolB string
}

func CreateConstantProductExchange(cped CpeDescriptor, tick uint64) ConstantProductExchange {
	var cpe ConstantProductExchange = ConstantProductExchange {
		ReserveA: tokens.CreateTokenReserve(tokens.Descriptor {
			Amount: cped.AmountA,
			Symbol: cped.SymbolA,
		}),
		ReserveB: tokens.CreateTokenReserve(tokens.Descriptor {
			Amount: cped.AmountB,
			Symbol: cped.SymbolB,
		}),
	}

	cpe.Auditer = candles.CreateAuditer(globals.DefaultAuditerBufferSize, cpe.SpotPriceA(), tick)
	return cpe
}

func (exg ConstantProductExchange) SpotPriceA() *uint256.Int {
	// NOTE the price of a in terms of b is calculated via
	// B / A
	var spot *uint256.Int = uint256.NewInt(1)
	spot.Div(exg.ReserveB.Amount, exg.ReserveA.Amount)
	return spot
}

func (exg ConstantProductExchange) SpotPriceB() *uint256.Int {
	var spot *uint256.Int = uint256.NewInt(1)
	spot.Div(exg.ReserveA.Amount, exg.ReserveB.Amount)
	return spot
}

func (cpe *ConstantProductExchange) Merge(feat ConstantProductExchange) {
	cpe.ReserveA = feat.ReserveA
	cpe.ReserveB = feat.ReserveB
	cpe.Auditer = feat.Auditer
}

func (cpe ConstantProductExchange) String() string {
	return fmt.Sprintf("{ reserveA: %v; reserveB: %v, audit: %v }", cpe.ReserveA, cpe.ReserveB, cpe.Auditer)
}

func (cpe ConstantProductExchange) Clone() ConstantProductExchange {
	return ConstantProductExchange {
		ReserveA: cpe.ReserveA.Clone(),
		ReserveB: cpe.ReserveB.Clone(),
		Auditer: cpe.Auditer.Clone(),
	}
}

func (cpe ConstantProductExchange) Hash() uint64 {
	return xxhash.Sum64([]byte(cpe.String()))
}
