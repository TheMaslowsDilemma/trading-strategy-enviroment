package strategy

import (
	"tse-p1/market"
	"tse-p1/candles"
)

type Strategy interface {
	Decide(candles []candles.Candle, i int) (market.Action, float64)
	GetName() string
}