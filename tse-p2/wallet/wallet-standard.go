package wallet

import (
	"tse-p2/ledger"
	"tse-p2/token"
)

func InitDefaultWallet(l *ledger.Ledger) ledger.LedgerAddr {
	rs := []token.TokenReserve {
        token.TokenReserve {
            Symbol: "usd",
            Amount: 10000.0,
        },
        token.TokenReserve {
            Symbol: "eth",
            Amount: 0.0,
        },
    }
    // Init wallet takes care of TokenReserve Init by way of AddReserve
    return InitWallet(rs, l)
}