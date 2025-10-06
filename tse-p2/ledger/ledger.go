package ledger

import (
    "fmt"
    "crypto/sha256"
)

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
