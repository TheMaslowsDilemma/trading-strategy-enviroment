package ledger

import (
    "fmt"
    "crypto/sha256"
    "math/rand"
)

type LedgerFetcher func(LedgerAddr) LedgerItem

type LedgerAddr uint64

type LedgerItem interface {
    Hash()      [sha256.Size]byte
    Copy()      LedgerItem
    String()    string
}

type Ledger map[LedgerAddr]LedgerItem

func (l Ledger) GetItemString(id LedgerAddr) (string, error) {
    var li LedgerItem = l[id]
    if li == nil {
        return "", fmt.Errorf("ledger contains no value for id: %v", id)
    }
    return li.String(), nil
}

func RandomLedgerAddr() LedgerAddr {
    return LedgerAddr(uint64(rand.Uint32()) << 32 | uint64(rand.Uint32()))
}

func Merge(main *Ledger, feature Ledger) uint  {
    var (
        laddr   LedgerAddr
        litem   LedgerItem
        ftcount uint
    )

    for laddr, litem = range feature {
        (*main)[laddr] = litem.Copy()
        ftcount += 1
    }
    return ftcount
}
