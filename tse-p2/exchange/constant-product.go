package exchange

import (
        "fmt"
        "crypto/sha256"
	"tse-p2/token"
        "tse-p2/ledger"
        "tsse-p2/candles"
)

type ConstantProductExchange struct {
	TokenReserveA	ledger.LedgerAddr
	TokenReserveB	ledger.LedgerAddr
        CandleAuditer   ledger.LedgerAddr
}

// Need to init against a ledger... so we can create token reserves too
func InitConstantProductExchange(syma, symb string, na, nb uint64, l ledger.Ledger) {
    var (
        exaddr  ledger.LedgerAddr
        ataddr  ledger.LedgerAddr
        btaddr  ledger.LedgerAddr
        caaddr  ledger.ledgerAddr
        tkra    token.TokenReserve
        tkrb    token.TokenReserve
        ca      candles.CandleAudit
        ex      ConstantProductExchange
    )

    exaddr = ledger.RandomLedgerAddr()
    ataddr = ledger.RandomLedgerAddr()
    btaddr = ledger.RandomLedgerAddr()
    caaddr = ledger.RandomLedgerAddr()

    tkra = token.NewTokenReserve(syma, na)
    tkrb = token.NewTokenReserve(symb, nb)
    ca = candles.NewCandleAuditer(10) // TODO hardcoded number to global or param
    ex = ConstantProductExchange {
        TokenReserveA: ataddr,
        TokenReserveB: btaddr,
        CandleAuditer: caaddr,
    }

    // place ledger items on the ledger //
    l[exaddr] = ex
    l[caaddr] = ca
    l[ataddr] = tkra
    l[btaddr] = tkrb
}

func (cpe ConstantProductExchange) Copy() ledger.LedgerItem {
    return ConstantProductExchange() {
        TokenReserveA: cpe.TokenReserveA,
        TokenReserveB: cpe.TokenReserveB,
        CandleAuditer: cpe.CandleAuditer.
    }
}

func (cpe ConstantProductExchange) String() string {
    return fmt.Sprintf(
        "{ rsv-a: %v, rsv-b: %v, cndla: %v }",
        cpe.TokenReserveA,
        cpe.TokenReserveB,
        cpe.CandleAuditer,
    )
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

    return (tkrA.Amount * tkrB.Amount) / (tkrA.Amount + ain), nil
}


func (cpe ConstantProductExchange) SwapBForA(l ledger.Ledger, bin uint64) (uint64, error) {
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

