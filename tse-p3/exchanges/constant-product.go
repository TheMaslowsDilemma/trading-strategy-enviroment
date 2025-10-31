package exchanges

import (
	"fmt"
	"tse-p3/tokens"
	"github.com/cespare/xxhash"
)

type ConstantProductExchange struct {
	ReserveA	tokens.TokenReserve
	ReserveB	tokens.TokenReserve
}

type CpeDescriptor struct {
	AmountA	uint64
	AmountB uint64
	SymbolA string
	SymbolB string
}

func CreateConstantProductExchange(cped CpeDescriptor) ConstantProductExchange {
	return ConstantProductExchange {
		ReserveA: tokens.CreateTokenReserve(cped.AmountA, cped.SymbolA),
		ReserveB: tokens.CreateTokenReserve(cped.AmountB, cped.SymbolB),
	}
}

func (cpe ConstantProductExchange) Merge(feat ConstantProductExchange) {
	(&cpe.ReserveA).Merge(feat.ReserveA)
	(&cpe.ReserveB).Merge(feat.ReserveB)
	
}

func (cpe ConstantProductExchange) String() string {
	return fmt.Sprintf("{ reserveA: %v; reserveB: %v }", cpe.ReserveA, cpe.ReserveB)
}

func (cpe ConstantProductExchange) Clone() ConstantProductExchange {
	return ConstantProductExchange {
		ReserveA: cpe.ReserveA.Clone(),
		ReserveB: cpe.ReserveB.Clone(),
	}
}

func (cpe ConstantProductExchange) Hash() uint64 {
	return xxhash.Sum64([]byte(cpe.String()))
}
