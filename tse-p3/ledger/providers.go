package ledger

import (
	"tse-p3/wallets"
)

type RateProvider func (sym, inTermsOf string) (float64, error)
type WalletProvider func (waddr Addr) (wallets.Wallet, error)


