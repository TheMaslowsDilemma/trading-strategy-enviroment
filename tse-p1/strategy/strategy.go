package strategy

import (
	"tse-p1/market"
	"tse-p1/candles"
)

type Strategy interface {
<<<<<<< Updated upstream
	Decide(candles []candles.Candle, currentIndex int) (market.Action, float64)
=======
	Decide(candles []candles.Candle, i int) (market.Action, float64)
>>>>>>> Stashed changes
	GetName() string
}