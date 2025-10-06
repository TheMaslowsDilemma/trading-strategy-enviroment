package ledger

import (
    "fmt"
    "crypto/sha256"
)

type LedgerItem interface {
    Hash()      [sha256.Size]byte
    Copy()      LedgerItem
    String()    string
}

type Ledger map[uint64]LedgerItem

func (l Ledger) GetItemString(id uint64) (string, error) {
    var li LedgerItem = l[id]
    if li == nil {
        return "", fmt.Errorf("ledger contains no value for id: %v", id)
    }
    return li.String(), nil
}
