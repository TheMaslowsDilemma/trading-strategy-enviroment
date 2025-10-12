package candles

import (
    "fmt"
)

type Candle struct {
    High        float64
    Low         float64
    Open        float64
    Close       float64
    Volume      float64
}

func (c *Candle) Start(price float64, volume float64) {
    c.Open = price
    c.Add(price, volume)
}

func (c *Candle) Add(price float64, volume float64) {
    if c.High < price {
        c.High = price
    }
    if c.Low > price {
        c.Low = price
    }
    c.Close = price
    c.Volume += volume
}

func (c Candle) String() string {
    return fmt.Sprintf("{ h: %.3f, l: %.3f, o: %.3f, c: %.3f, v: %.3v }",
        c.High,
        c.Low,
        c.Open,
        c.Close,
        c.Volume,
    )
}


func (c Candle) Copy() Candle {
    return Candle {
        High: c.High,
        Low: c.Low,
        Open: c.Open,
        Close: c.Close,
        Volume: c.Volume,
    }
}
