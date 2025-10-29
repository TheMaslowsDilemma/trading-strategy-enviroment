package ledger

import (
	"fmt"
	"tse-p3/exchange"
	"tse-p3/wallet"
)

type Ledger struct {
	Wallets		[Addr]wallets.Wallet
	Exchanges	[Addr]exchanges.Exchange
}


func CreateLedger() Ledger {
	var (
		ws [Addr]wallets.Wallet
		es [Addr]exchanges.Exchange
	)

	ws = make([Addr]wallets.Wallet)
	es = make([Addr]exchanges.Exchange)

	return Ledger {
		Wallets: ws,
		Exchanges: es,
	}
}



