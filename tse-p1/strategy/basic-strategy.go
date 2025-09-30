package strategy

import (
	"tse-p1/candles"
	"tse-p1/market"
)

type SimpleMAStrategy struct {
	ShortPeriod int
	LongPeriod  int
}

/**
 *  returns an action {buy, sell, hold} and
 * 	a confidence [0,1] which is basically how much
**/
func (s *SimpleMAStrategy) Decide(candles []candles.Candle, currentIndex int) (market.Action, float64) {
	if currentIndex < s.LongPeriod {
		return market.Hold, 1 // Not enough data yet
	}

	shortMA := calculateMA(candles, currentIndex, s.ShortPeriod)
	longMA := calculateMA(candles, currentIndex, s.LongPeriod)

	prevShortMA := calculateMA(candles, currentIndex-1, s.ShortPeriod)
	prevLongMA := calculateMA(candles, currentIndex-1, s.LongPeriod)


	if prevShortMA <= prevLongMA && shortMA > longMA {
		return market.Sell, 1 // Crossover up
	}
	if prevShortMA >= prevLongMA && shortMA < longMA {
		return market.Buy, 1 // Crossover down
	}
	return market.Hold, 1
}

func calculateMA(candles []candles.Candle, currentIndex, period int) float64 {
	sum := 0.0
	start := currentIndex - period + 1
	for i := start; i <= currentIndex; i++ {
		sum += candles[i].Close
	}
	return sum / float64(period)
}
