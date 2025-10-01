package strategy

import (
	"math"
	"tse-p1/candles"
	"tse-p1/market"
)

type SimpleStrategy struct {
	ShortInterval int
	LongInterval int
}

func (s *SimpleStrategy) Decide(cs []candles.Candle, i int) (market.Action, float64) {
	var (
		ss float64 // short slope //
		ls float64 // long slope //
	)

	// verify we have seen enough candles
	if i < s.LongInterval {
		return market.Hold, 1
	}

	// this could happen technically - so make sure it cant
	if s.ShortInterval >= s.LongInterval {
		return market.Hold, 1
	}

	// Find Slopes //
	ss = linearRegressionSlope(cs, i - s.ShortInterval, i)
	ls = linearRegressionSlope(cs, i - s.LongInterval, i - s.ShortInterval)

	if ss >= 0 {
		if ls >= 0 {
			return market.Buy, Sigmoid(ss + ls)
		} else {
			return market.Sell, Sigmoid(ss)
		}
	} else {
		if ls >= 0 {
			return market.Hold, 1
		} else {
			return market.Sell, Sigmoid(math.Abs(ss))
		}
	}

	return market.Hold, 1
}

func (s *SimpleStrategy) GetName() string {
	return "SimpleStrategy"
}

func linearRegressionSlope(cs []candles.Candle, startIdx int, endIdx int) float64 {
    n := float64(endIdx - startIdx)
    if n <= 1 {
        return 0
    }

    var sumX, sumY, sumXY, sumXX float64
    for i := startIdx; i < endIdx; i++ {
        x := float64(i - startIdx)
        y := cs[i].High // Use high prices for trend
        sumX += x
        sumY += y
        sumXY += x * y
        sumXX += x * x
    }

    return (n*sumXY - sumX*sumY) / (n*sumXX - sumX*sumX)
}

func Sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

func calculateSlope(a float64, b float64, interval float64) float64 {
	return (a - b) / interval
}

