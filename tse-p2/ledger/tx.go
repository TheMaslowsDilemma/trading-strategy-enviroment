package ledger

type TxResult uint8
const (
    TxFail = iota
    TxPass
)

type Tx interface {
    Apply(tick uint64, l Ledger) (Ledger, error)
    Notify(result TxResult)
}
