package wallet

import (
	"fmt"
        "tse-p2/token"
        "tse-p2/ledger"
	"crypto/sha256"
)

type Wallet struct {
	Reserves	[]ledger.LedgerAddr
}

func InitWallet(rs []token.TokenReserve, l *ledger.Ledger) ledger.LedgerAddr {
    var (
        wlt     Wallet
        wltaddr ledger.LedgerAddr
        err     error
    )
    
    wlt = Wallet{ Reserves: make([]ledger.LedgerAddr, 0) }
    wltaddr = ledger.RandomLedgerAddr()
    
    for _, r := range rs {
        fmt.Printf("adding %v with %v amt\n", r.Symbol, r.Amount)
        err = (&wlt).AddReserve(r.Symbol, r.Amount, l)
        if err != nil {
            fmt.Printf("err adding \"%v\" to wallet: %v", r.Symbol, err)
        }
    }

    fmt.Println(wlt)
    (*l)[wltaddr] = wlt
    return wltaddr
}

func (w *Wallet) AddReserve(sym string, amt float64, l *ledger.Ledger) error {
    var tkaddr ledger.LedgerAddr
    if _, err := w.GetReserveAddr(sym, *l); err == nil {
        return fmt.Errorf("wallet already contains reserve %v", sym)
    }

    tkaddr = token.InitTokenReserve(sym, amt, l)
    w.Reserves = append(w.Reserves, tkaddr)

    return nil
}

func (w Wallet) String() string {
	return fmt.Sprintf("reserves: %v", w.Reserves)
}

func (w Wallet) Hash() [sha256.Size]byte {
	return sha256.Sum256([]byte(w.String()))
}

func (w Wallet) Copy() ledger.LedgerItem {
    var rs []ledger.LedgerAddr = make([]ledger.LedgerAddr, len(w.Reserves))

    for i, r := range w.Reserves {
        rs[i] = r
    }

    return Wallet {
        Reserves: rs,
    }
}

func (w Wallet) ContainsReserve(raddr ledger.LedgerAddr) bool {
    var i int
    for i = 0; i < len(w.Reserves); i++ {
        if w.Reserves[i] == raddr {
            return true
        }
    }
    return false
}

func (w Wallet) GetReserveAddr(sym string, l ledger.Ledger) (ledger.LedgerAddr, error) {
    var (
        i       int
        tkr     *token.TokenReserve
        err     error
    )

    for i = 0; i < len(w.Reserves); i++ {
        tkr, err = token.TkrFromLedgerItem(l[w.Reserves[i]])
        if err != nil {
            continue; // we can ignore these errors
        }
        if tkr.Symbol == sym {
            return w.Reserves[i], nil
        }
    }
    return 0, fmt.Errorf("wallet reserve for \"%v\" DNE.", sym)
}

func WalletFromLedgerItem(li ledger.LedgerItem) (*Wallet, error) {
     var (
         wlt    Wallet
         ok     bool
     )

     if li == nil {
         return nil, fmt.Errorf("cannot cast wallet from nil ledger item.")
     }

     if wlt, ok = li.(Wallet); ok {
         return &wlt, nil
     }

     return nil, fmt.Errorf("cannot cast wallet from non-wallet ledger item.")
}
