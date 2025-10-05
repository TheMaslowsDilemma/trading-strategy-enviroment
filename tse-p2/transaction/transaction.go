package transaction

import (
    "tse-p2/ledger"
)

type Tx interface {
    Apply(l *ledger.Ledger) *ledger.Ledger, error
}
