package market

import (
	"tse-p1/candles"
)

type Action uint8
const (
	Hold Action = iota
	Buy
	Sell
)

type Strategy interface {
	Decide(candles []candles.Candle, idx uint)	
}
