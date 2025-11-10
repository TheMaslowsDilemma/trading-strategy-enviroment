package strategies

import (
	"tse-p3/candles"
)

type Action uint8
const (
	Buy = iota
	Sell
	Hold
)

type Strategy interface {
	Decide(cs []candles.Candle)	(Action, float64)
}