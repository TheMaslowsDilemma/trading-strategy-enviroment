package ledger

type Transaction interface {
    Apply(l Ledger)
}
