package ledger

import (
    "fmt"
)

type LedgerItem interface {
    Hash()      []byte
    Copy()      LedgerItem
    String()    string
}

type Ledger map[uint64]LedgerItem

func (l *Ledger) GetItemString(id uint64) (string, err) {
    var li LedgerItem = l[id]
    if li == nil {
        return "", fmt.Errorf("ledger contains no value for id: %v", id)
    }
    return li.String(), nil
}
