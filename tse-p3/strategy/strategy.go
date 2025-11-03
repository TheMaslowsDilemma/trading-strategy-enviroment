package strategies


type Action uint8
const (
	Buy = iota
	Sell
	Hold
)

type Strategy interface {
	Decide()	Action
}