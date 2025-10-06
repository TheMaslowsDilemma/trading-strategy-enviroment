package exchange

import (
        "fmt"
        "crypto/sha256"
	"tse-p2/ledger"
)

type ConstantProductExchange struct {
	TokenReserveA	ledger.LedgerAddr
	TokenReserveB	ledger.LedgerAddr
}

func (cpe ConstantProductExchange) Copy() ledger.LedgerItem {
    return ConstantProductExchange() {
        TokenReserveA: cpe.TokenReserveA,
        TokenReserveB: cpe.TokenReserveB,
    }
}

func (cpe ConstantProductExchange) String() string {
    return fmt.Sprintf("{ rsv-a: %v, rsv-b: %v}", cpe.TokenReserveA, cpe.TokenReservB)
}

func (cpe ConstantProductExchange) Hash() [sha256.Size]byte {
    return sha256.Sum256([]byte(cpe.String()))
}

func CpeFromLedgerItem(li *ledger.LedgerItem) *ConstantProductExchange, error {
    var (
        cpe     *ConstantProductExchange
        ok      bool
    )

    if li == nil {
        return nil, fmt.Errorf("cannot cast cpe from nil ledger item.")
    }
    
    if cpe, ok = li.(*ConstantProductExchange); ok {
        return cpe, nil
    }
    return nil, fmt.Errorf("cannot cast cpe from non-cpe ledger item.")
}


