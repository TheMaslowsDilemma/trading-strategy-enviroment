package txs

import (
	"tse-p3/ledger"
)

type TxResult uint8
const (
	TxFail = iota
	TxPass
)

type Tx interface {
	Apply(tick uint64, l *ledger.Ledger) (ledger.Ledger, error)
	Notify(result TxResult)
}
