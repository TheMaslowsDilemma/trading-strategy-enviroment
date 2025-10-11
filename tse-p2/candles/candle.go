package candles

import (
    "fmt"
)

type Candle struct {
    High        float64
    Low         float64
    Open        float64
    Close       float64
    Volume      uint64
}

func (c *Candle) Open(cost float64, volume uint64) {
    c.Open = cost
    c.Add(cost, volume)
}

func (c *Candle) Add(cost float64, volume uint64) {
    if c.High < cost {
        c.High = cost
    }
    if c.Low > cost {
        c.Low = cost
    }
    c.Close = cost
    c.Volume += volume
}

func (c Candle) String() {
    return fmt.Sprintf("{ h: %.3f, l: %.3f, o: %.2f, c: %.2f, v: %v }",
        c.High,
        c.Low,
        c.Open,
        c.Close,
        c.Volume
    }
}

func (c Candle) Copy() {
    return Candle {
        High: c.High,
        Low: c.Low,
        Open: c.Open,
        Close: c.Close,
        Volume: c.Volume
    }
}
