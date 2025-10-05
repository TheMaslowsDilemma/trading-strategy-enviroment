package ledger

type LedgerItem interface {
    Hash() []byte
    Copy() LedgerItem
}

type Ledger map[uint64]LedgerItem
