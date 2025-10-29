package ledger

type LedgerItem interface {
	Address() Addr
	Hash() uint64
	String() string
}

