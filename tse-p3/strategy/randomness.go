package strategies

import (
	"math/rand"
	"tse-p3/candles"
)

type RandomStrategy struct { }

func (rs RandomStrategy) Decide(cs []candles.Candle) (Action, float64) {
	var act float64 = rand.Float64()
	var cnf float64 = rand.Float64()

	if act > 0.66 {
		return Buy, cnf
	} else if act > 0.33 {
		return Sell, cnf
	} else {
		return Hold, cnf
	}
}
