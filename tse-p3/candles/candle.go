package candles

import (
	"fmt"
	"github.com/holiman/uint256"
)

type Candle struct {
	Ts	uint64 // either sim-tick or timestamp
	High	*uint256.Int
	Open	*uint256.Int
	Low	*uint256.Int
	Close	*uint256.Int
}

func CreateCandle(price *uint256.Int, ts uint64) Candle {
	return Candle {
		Open: price,
		High: price,
		Low : price,
		Close: price,
		Ts: ts,
	}
}

func (c *Candle) Start(price *uint256.Int, ts uint64) {
	c.Open = price
	c.Low = price
	c.High = price
	c.Close = price
	c.Ts = ts
}

func (c *Candle) Add(price *uint256.Int) {
	if c.High.Lt(price) {
		c.High = price
	}
	if c.Low.Gt(price) {
		c.Low = price
	}
	c.Close = price
}

func (c Candle) String() string {
	return fmt.Sprintf("{o: %v, h: %v, l: %v, c: %v }", c.Open, c.High, c.Low, c.Close)
}

func (c Candle) Clone() Candle {
	return Candle {
		Ts: c.Ts,
		Open: c.Open.Clone(),
		High: c.High.Clone(),
		Low:  c.Low.Clone(),
		Close: c.Close.Clone(),
	}
}

// NOTE the way we merge is somewhat arbitrary
func (c *Candle) Merge(feat Candle) {
	if c.Ts != feat.Ts {
		// NOTE we might want some warning / error here...
		return
	}

	// This part is arbitrary
	if c.Open == nil {
		c.Open = feat.Open
	}
	if c.Close == nil {
		c.Close = feat.Close
	}

	if c.High.Lt(feat.High) {
		c.High = feat.High
	}
	if c.Low.Gt(feat.Low) {
		c.Low = feat.Low
	}
}
