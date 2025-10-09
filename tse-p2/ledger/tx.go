package ledger

type Tx interface {
    Apply(l Ledger) (Ledger, error)
}
