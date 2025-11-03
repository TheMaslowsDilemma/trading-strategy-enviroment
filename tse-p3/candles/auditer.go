package candles

import (
	"fmt"
	"github.com/holiman/uint256"
)

type Auditer struct {
	CandleBuffer	CandleBuffer
	ActiveCandle	Candle
}

// NOTE an auditer should be initialized with a starting price and timestamp -- avoids bad candle at start.
func CreateAuditer(buffsize uint, initPrice *uint256.Int, ts uint64) Auditer {
	return Auditer {
		CandleBuffer: CreateCandleBuffer(buffsize),
		ActiveCandle: CreateCandle(initPrice, ts),
	}
}

func (a *Auditer) Audit(price *uint256.Int, tick uint64) {
	if tick != a.ActiveCandle.Ts {
		a.CandleBuffer.Push(a.ActiveCandle)
		(&a.ActiveCandle).Start(price, tick)
	} else {
		(&a.ActiveCandle).Add(price)
	}
}

// returns list of candles in order of ts ascending
func (a Auditer) GetCandles() []Candle {
	var cs []Candle = a.CandleBuffer.GetCandles()
	cs = append(cs, a.ActiveCandle)
	return cs
}

func (a *Auditer) String() string {
	return fmt.Sprintf("{ active: %v }", a.ActiveCandle)
}

func (a *Auditer) Clone() *Auditer {
	return &Auditer {
		CandleBuffer: a.CandleBuffer.Clone(),
		ActiveCandle: a.ActiveCandle.Clone(),
	}
}
