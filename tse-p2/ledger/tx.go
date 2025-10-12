package ledger

type Tx interface {
    Apply(tick uint64, l Ledger) (Ledger, error)
}
