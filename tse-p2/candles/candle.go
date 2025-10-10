package candles

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
    c.Volume += volume
}

func (c *Candle) Close(cost float64, volume uint64) {
    c.Close = cost
    c.Add(cost, volume)
}
