package strategy

import (
	"math"
	"tse-p2/candles"
)

// MomentumStrategy: Decides based on price momentum over a lookback period.
// Positive momentum suggests buy, negative suggests sell, with threshold for hold.
type MomentumStrategy struct {
	Lookback int     // Periods to calculate momentum over
	Threshold float64 // Minimum absolute momentum to trigger buy/sell
}

func (s MomentumStrategy) Decide(cs []candles.Candle) (Action, float64) {
	i := len(cs)
	if i < s.Lookback+1 {
		return Hold, 1
	}

	currentPrice := cs[i-1].Close
	pastPrice := cs[i-1-s.Lookback].Close
	momentum := (currentPrice - pastPrice) / pastPrice

	absMomentum := math.Abs(momentum)
	if absMomentum < s.Threshold {
		return Hold, 1
	}

	confidence := Sigmoid(absMomentum / s.Threshold) * 0.5
	if momentum > 0 {
		return Buy, confidence
	}
	return Sell, confidence
}

func (s MomentumStrategy) GetName() string {
	return "MomentumStrategy"
}

// VolatilityBreakoutStrategy: Trades on volatility breakouts using ATR (Average True Range).
// Buys on upside breakout, sells on downside.
type VolatilityBreakoutStrategy struct {
	ATRPeriod int     // Period for ATR calculation
	Multiplier float64 // Multiplier for breakout threshold
}

func (s VolatilityBreakoutStrategy) Decide(cs []candles.Candle) (Action, float64) {
	i := len(cs)
	if i < s.ATRPeriod+1 {
		return Hold, 1
	}

	atr := calculateATR(cs, i, s.ATRPeriod)
	prevClose := cs[i-2].Close
	currentHigh := cs[i-1].High
	currentLow := cs[i-1].Low

	upBreak := currentHigh > prevClose + atr*s.Multiplier
	downBreak := currentLow < prevClose - atr*s.Multiplier

	if upBreak && !downBreak {
		return Buy, Sigmoid(s.Multiplier) * 0.5
	} else if downBreak && !upBreak {
		return Sell, Sigmoid(s.Multiplier) * 0.5
	}
	return Hold, 1
}

func (s VolatilityBreakoutStrategy) GetName() string {
	return "VolatilityBreakoutStrategy"
}

func calculateATR(cs []candles.Candle, endIdx int, period int) float64 {
	var sum float64
	for i := endIdx - period; i < endIdx; i++ {
		hl := cs[i].High - cs[i].Low
		hpc := math.Abs(cs[i].High - cs[i-1].Close)
		lpc := math.Abs(cs[i].Low - cs[i-1].Close)
		tr := math.Max(hl, math.Max(hpc, lpc))
		sum += tr
	}
	return sum / float64(period)
}

// MeanReversionStrategy: Assumes prices revert to mean; buys below mean, sells above.
// Uses simple moving average as mean.
type MeanReversionStrategy struct {
	SMAPeriod int     // Period for simple moving average
	Deviation float64 // Deviation threshold as percentage
}

func (s MeanReversionStrategy) Decide(cs []candles.Candle) (Action, float64) {
	i := len(cs)
	if i < s.SMAPeriod {
		return Hold, 1
	}

	sma := calculateSMA(cs, i, s.SMAPeriod)
	currentPrice := cs[i-1].Close
	dev := math.Abs(currentPrice - sma) / sma

	if dev < s.Deviation {
		return Hold, 1
	}

	confidence := Sigmoid(dev / s.Deviation) * 0.5
	if currentPrice < sma {
		return Buy, confidence
	}
	return Sell, confidence
}

func (s MeanReversionStrategy) GetName() string {
	return "MeanReversionStrategy"
}

func calculateSMA(cs []candles.Candle, endIdx int, period int) float64 {
	var sum float64
	for i := endIdx - period; i < endIdx; i++ {
		sum += cs[i].Close
	}
	return sum / float64(period)
}

// RandomWalkStrategy: For baseline diversity, makes random decisions with controlled probability.
// Useful for simulating noise traders.
type RandomWalkStrategy struct {
	BuyProb  float64 // Probability to buy
	SellProb float64 // Probability to sell (hold if neither)
}

func (s RandomWalkStrategy) Decide(cs []candles.Candle) (Action, float64) {
	if len(cs) < 2 {
		return Hold, 1
	}

	randVal := math.Mod(cs[len(cs)-1].Close*100, 1.0) // Pseudo-random based on last close for determinism in sim
	if randVal < s.BuyProb {
		return Buy, s.BuyProb
	} else if randVal < s.BuyProb+s.SellProb {
		return Sell, s.SellProb
	}
	return Hold, 1
}

func (s RandomWalkStrategy) GetName() string {
	return "RandomWalkStrategy"
}
