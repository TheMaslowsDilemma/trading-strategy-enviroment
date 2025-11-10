package candles

import (
	"fmt"
)

type Candle struct {
	Ts		uint64
	High	float64
	Open	float64
	Low		float64
	Close	float64
}

func CreateCandle(price float64, ts uint64) Candle {
	return Candle {
		Open: price,
		High: price,
		Low : price,
		Close: price,
		Ts: ts,
	}
}

func (c *Candle) Start(price float64, ts uint64) {
	c.Open = price
	c.Low = price
	c.High = price
	c.Close = price
	c.Ts = ts
}

func (c *Candle) Add(price float64) {
	if c.High >= price {
		c.High = price
	}
	if c.Low <= price {
		c.Low = price
	}
	c.Close = price
}

func (c Candle) String() string {
	return fmt.Sprintf("{ o: %v, h: %v, l: %v, c: %v }", c.Open, c.High, c.Low, c.Close)
}

func (c Candle) Clone() Candle {
	return Candle {
		Ts: c.Ts,
		Open: c.Open,
		High: c.High,
		Low:  c.Low,
		Close: c.Close,
	}
}

// NOTE the way we merge is somewhat arbitrary
func (c *Candle) Merge(feat Candle) {
	if c.Ts != feat.Ts {
		// NOTE we might want some warning / error here...
		return
	}

	// This part is arbitrary
	if c.Open == 0 {
		c.Open = feat.Open
	}
	if c.Close == 0 {
		c.Close = feat.Close
	}

	if c.High < feat.High {
		c.High = feat.High
	}
	if c.Low > feat.Low {
		c.Low = feat.Low
	}
}
