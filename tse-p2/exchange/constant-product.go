package exchange

import (
	"fmt"
	"crypto/sha256"
	"tse-p2/token"
	"tse-p2/ledger"
	"tse-p2/candles"
)

type ConstantProductExchange struct {
    TkrAddrA    ledger.LedgerAddr
    TkrAddrB    ledger.LedgerAddr
    CndlAddr    ledger.LedgerAddr // ledger addr of CandleAudit
}

// Need to init against a ledger... so we can create token reserves too
func InitConstantProductExchange(symbA, symbB string, cntA, cntB float64, l *ledger.Ledger) ledger.LedgerAddr {
    var (
        exaddr  ledger.LedgerAddr
        ex      ConstantProductExchange
    )

    exaddr = ledger.RandomLedgerAddr()
    ex = ConstantProductExchange {
        TkrAddrA: token.InitTokenReserve(symbA, cntA, l),
        TkrAddrB: token.InitTokenReserve(symbB, cntB, l),
        CndlAddr: candles.InitCandleAudit(10, l), // NOTE hard coded
    }

    (*l)[exaddr] = ex
    return exaddr
}

func (cpe ConstantProductExchange) Copy() ledger.LedgerItem {
    return ConstantProductExchange {
        TkrAddrA: cpe.TkrAddrA,
        TkrAddrB: cpe.TkrAddrB,
        CndlAddr: cpe.CndlAddr,
    }
}

func (cpe ConstantProductExchange) String() string {
    return fmt.Sprintf(
        "{ rsv-a: %v, rsv-b: %v, cndla: %v }",
        cpe.TkrAddrA,
        cpe.TkrAddrB,
        cpe.CndlAddr,
    )
}

func (cpe ConstantProductExchange) Hash() [sha256.Size]byte {
    return sha256.Sum256([]byte(cpe.String()))
}

func (cpe ConstantProductExchange) GetPriceA(l ledger.Ledger) (float64, error) {
    tkrA, err := token.TkrFromLedgerItem(l[cpe.TkrAddrA])
    if err != nil {
        return 0, fmt.Errorf("tkrA failed: %v", err)
    }
    tkrB, err := token.TkrFromLedgerItem(l[cpe.TkrAddrB])
    if err != nil {
        return 0, fmt.Errorf("tkrB failed: %v", err)
    }
    if tkrA.Amount == 0 {
        return 0, fmt.Errorf("token A reserve is zero")
    }
    return tkrB.Amount / tkrA.Amount, nil
}

func (cpe ConstantProductExchange) GetPriceB(l ledger.Ledger) (float64, error) {
    var (
        tkra    *token.TokenReserve
        tkrb    *token.TokenReserve
        err     error
    )

    tkrb, err = token.TkrFromLedgerItem(l[cpe.TkrAddrB])
    if err != nil {
        return 0, fmt.Errorf("failed to cast tkr b: %v", err)
    }

    tkra, err = token.TkrFromLedgerItem(l[cpe.TkrAddrA])
    if err != nil {
        return 0, fmt.Errorf("failed to cast tkr a: %v", err)
    }

    if tkra.Amount == 0 {
        return 0, fmt.Errorf("tkr a amount is zero")
    }
    return tkra.Amount / tkrb.Amount, nil
}

func (cpe ConstantProductExchange) SwapAForB(l ledger.Ledger, ain float64) (float64, error) {
    var (
        tkrA    *token.TokenReserve
        tkrB    *token.TokenReserve
        err     error
    )
    
    tkrA, err = token.TkrFromLedgerItem(l[cpe.TkrAddrA])
    if err != nil {
        return 0, fmt.Errorf("tkrA failed: %v", err)
    }

    tkrB, err = token.TkrFromLedgerItem(l[cpe.TkrAddrB])
    if err != nil {
        return 0, fmt.Errorf("tkrB failed: %v", err)
    }

    return (tkrA.Amount * tkrB.Amount) / (tkrA.Amount + ain), nil
}


func (cpe ConstantProductExchange) SwapBForA(l ledger.Ledger, bin float64) (float64, error) {
    var (
        tkrA    *token.TokenReserve
        tkrB    *token.TokenReserve
        err     error
    )
    
    tkrA, err = token.TkrFromLedgerItem(l[cpe.TkrAddrA])
    if err != nil {
        return 0, fmt.Errorf("tkrA failed: %v", err)
    }

    tkrB, err = token.TkrFromLedgerItem(l[cpe.TkrAddrB])
    if err != nil {
        return 0, fmt.Errorf("tkrB failed: %v", err)
    }

    // k = B * A // k is constant
    return (tkrB.Amount * tkrA.Amount) / (tkrB.Amount + bin), nil
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
