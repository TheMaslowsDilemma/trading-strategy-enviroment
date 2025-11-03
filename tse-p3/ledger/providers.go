package ledger

import (
	"tse-p3/wallets"
	"github.com/holiman/uint256"
)

type RateProvider func (sym, inTermsOf string) (*uint256.Int, error)
type WalletProvider func (waddr Addr) (wallets.Wallet, error)


