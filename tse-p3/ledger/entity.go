package ledger

type EntityType uint8

const (
	Wallet_t = iota
	Exchange_t
)

func (et EntityType) String() string {
	switch et {
	case Wallet_t:
		return "wallet"
	case Exchange_t:
		return "exchange"
	default:
		return "Unknown"
	}
}