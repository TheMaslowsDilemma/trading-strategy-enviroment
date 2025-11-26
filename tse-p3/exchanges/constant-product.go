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
	Name	string
	AmountA	uint64
	AmountB uint64
	SymbolA string
	SymbolB string
}

func CreateConstantProductExchange(cd CpeDescriptor, tick uint64) ConstantProductExchange {
	var cpe ConstantProductExchange = ConstantProductExchange {
		ReserveA: tokens.CreateTokenReserve(tokens.Descriptor {
			Amount: cd.AmountA,
			Symbol: cd.SymbolA,
		}),
		ReserveB: tokens.CreateTokenReserve(tokens.Descriptor {
			Amount: cd.AmountB,
			Symbol: cd.SymbolB,
		}),
	}

	cpe.Auditer = candles.CreateAuditer(globals.DefaultAuditerBufferSize, cpe.SpotPriceA(), tick)
	return cpe
}

func (exg ConstantProductExchange) SpotPriceA() float64 {
	// NOTE the price of a in terms of b is calculated via
	// B / A
	var (
		scaled *uint256.Int
		descld float64
	)

	scaled = uint256.NewInt(0)
	scaled.Div(scaled.Mul(exg.ReserveB.Amount, globals.TokenScaleFactor), exg.ReserveA.Amount)
	descld = scaled.Float64() / globals.TokenScaleFactorf64
	return descld
}

func (exg ConstantProductExchange) SpotPriceB() float64 {
	var (
		scaled *uint256.Int
		descld float64
	)

	scaled = uint256.NewInt(0)
	scaled.Div(scaled.Mul(exg.ReserveA.Amount, globals.TokenScaleFactor), exg.ReserveB.Amount)
	descld = scaled.Float64() / globals.TokenScaleFactorf64
	return descld
}

func (exg ConstantProductExchange) SwapAForB(amt_in *uint256.Int) *uint256.Int {
	var fiz *uint256.Int = uint256.NewInt(0)
	var	buz	*uint256.Int = uint256.NewInt(0)
	fiz.Mul(exg.ReserveB.Amount, amt_in)	// (B0 * amt_in)
	buz.Add(exg.ReserveA.Amount, amt_in)	// (A0 + amt_in)
	return fiz.Div(fiz, buz)				// (A0 * amt_in) / (B0 + amt_in)
}

func (exg ConstantProductExchange) SwapBForA(amt_in *uint256.Int) *uint256.Int {
	var fiz *uint256.Int = uint256.NewInt(0)
	var	buz	*uint256.Int = uint256.NewInt(0)

	// K0 = A0 * B0
	// B1 = (B0 + amt_in)
	// A1 = (A0 - amt_out)
	// A0 * B0 = A1 * B1
	// A0 * B0 = (A0 - amt_out) * (B0 + amt_in)
	// A0 * B0 / (B0 + amt_in) = (A0 - amt_out)

	// amt_out = A0 - A0 * B0 / (B0 + amt_in) 

	/** -- NOTE -- 
	 * the bove equation is correct-ish, but it leads to loss due to floor div.
	 * But the integer division means the trader will get more and exchange less
	**/

	// amt_out * (B0 + amt_in) = A0 * (B0 + amt_in) - (A0 * B0)
	// amt_out * (B0 + amt_in) = (A0 * B0) + (A0 * amt_in) - (A0 * B0)
	// amt_out = (A0 * amt_in) / (B0 + amt_in)

	fiz.Mul(exg.ReserveA.Amount, amt_in)	// (A0 * amt_in)
	buz.Add(exg.ReserveB.Amount, amt_in)	// (B0 + amt_in)
	return fiz.Div(fiz, buz)				// (A0 * amt_in) / (B0 + amt_in)
}

func (cpe *ConstantProductExchange) Merge(feat ConstantProductExchange) {
	cpe.ReserveA = feat.ReserveA
	cpe.ReserveB = feat.ReserveB
	cpe.Auditer = feat.Auditer
}

func (cpe ConstantProductExchange) String() string {
	return fmt.Sprintf("{ reserveA: %v; reserveB: %v, audit: %v}", cpe.ReserveA, cpe.ReserveB, cpe.Auditer)
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
