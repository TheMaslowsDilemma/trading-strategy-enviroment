package exchange

import (
        "fmt"
        "crypto/sha256"
	"tse-p2/token"
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

func (cpe ConstantProductExchange) SwapAForB(l ledger.Ledger, ain uint64) (uint64, error) {
    var (
        tkrA    token.TokenReserve
        tkrB    token.TokenReserve
        err     error
    )
    
    tkrA, err = token.TkrFromLedgerItem(l[cpe.TokenReserveA])
    if err != nil {
        return 0, fmt.Errorf("tkrA failed: %v", err)
    }

    tkrB, err = token.TkrFromLedgerItem(l[cpe.TokenReserveB])
    if err != nil {
        return 0, fmt.Errorf("tkrB failed: %v", err)
    }

    return (tkrA.Amount * tkrB.Amount) / (tkrA.Amount + ain), 0
}


func (cpe ConstantProductExchange) SwapBForA(l ledger.Ledger, bin uint64) uint64, error {
    var (
        tkrA    token.TokenReserve
        tkrB    token.TokenReserve
        err     error
    )
    
    tkrA, err = token.TkrFromLedgerItem(l[cpe.TokenReserveA])
    if err != nil {
        return 0, fmt.Errorf("tkrA failed: %v", err)
    }

    tkrB, err = token.TkrFromLedgerItem(l[cpe.TokenReserveB])
    if err != nil {
        return 0, fmt.Errorf("tkrB failed: %v", err)
    }

    // k = B * A // k is constant
    return (tkrB.Amount * tkrA.Amount) / (tkrB.Amount + bin), 0
}

func CpeFromLedgerItem(li ledger.LedgerItem) (*ConstantProductExchange, error) {
    var (
        cpe     ConstantProductExchange
        ok      bool
    )

    if li == nil {
        return nil, fmt.Errorf("cannot cast cpe from nil ledger item.")
    }
    
    if cpe, ok = li.(ConstantProductExchange); ok {
        return &cpe, nil
    }
    return nil, fmt.Errorf("cannot cast cpe from non-cpe ledger item.")
}


