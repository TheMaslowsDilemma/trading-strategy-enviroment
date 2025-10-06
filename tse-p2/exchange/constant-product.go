package exchange

import (
	"tse-p2/ledger"
)

type ConstantProductExchange struct {
	TokenReserveA	ledger.LedgerAddr
	TokenReserveB	ledger.LedgerAddr
}

