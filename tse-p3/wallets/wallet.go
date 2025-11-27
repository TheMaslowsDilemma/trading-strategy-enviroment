package wallets

import (
	"fmt"
	"tse-p3/tokens"
	"github.com/cespare/xxhash"
)

type Wallet struct {
	Name	string
	Reserve	tokens.TokenReserve
}

type WalletDescriptor struct {
	Name	string
	Amount	uint64
	Symbol	string
}

func CreateWallet(wd WalletDescriptor) Wallet {
	return Wallet {
		Name: wd.Name,
		Reserve: tokens.CreateTokenReserve(tokens.Descriptor {
			Amount: wd.Amount,
			Symbol: wd.Symbol,
		}),
	}
}

func (wlt Wallet) Merge(feat Wallet) {
	(&wlt.Reserve).Merge(wlt.Reserve)
}

func (wlt Wallet) Clone() Wallet  {
	return Wallet {
		Name: wlt.Name,
		Reserve: wlt.Reserve.Clone(),
	}
}

func (wlt Wallet) String() string {
	return fmt.Sprintf("{name: %v reserve: %v }", wlt.Name, wlt.Reserve)
}

func (wlt Wallet) Hash() uint64 {
	return xxhash.Sum64([]byte(wlt.String()))
}
