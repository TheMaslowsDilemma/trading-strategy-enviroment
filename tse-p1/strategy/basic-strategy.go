package strategy

import (
	"tse-p1/candles"
	"tse-p1/market"
)

type SimpleMAStrategy struct {
	ShortPeriod int
	LongPeriod  int
}

func (s *SimpleMAStrategy) Decide(candles []candles.Candle, currentIndex int) market.Action {
	if currentIndex < s.LongPeriod {
		return market.Hold // Not enough data yet
	}

	shortMA := calculateMA(candles, currentIndex, s.ShortPeriod)
	longMA := calculateMA(candles, currentIndex, s.LongPeriod)

	prevShortMA := calculateMA(candles, currentIndex-1, s.ShortPeriod)
	prevLongMA := calculateMA(candles, currentIndex-1, s.LongPeriod)

	if prevShortMA <= prevLongMA && shortMA > longMA {
		return market.Buy // Crossover up
	}
	if prevShortMA >= prevLongMA && shortMA < longMA {
		return market.Sell // Crossover down
	}
	return market.Hold
}

func calculateMA(candles []candles.Candle, currentIndex, period int) float64 {
	sum := 0.0
	start := currentIndex - period + 1
	for i := start; i <= currentIndex; i++ {
		sum += candles[i].Close
	}
	return sum / float64(period)
}
