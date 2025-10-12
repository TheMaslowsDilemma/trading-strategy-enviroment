package candles

import (
    "fmt"
)

type CandleCirq struct {
    Cs  []Candle
    Cap uint32
    Fnt uint32
    Cnt uint32
}

func NewCandleCirq(n uint32) CandleCirq {
    var cs []Candle = make([]Candle, int(n))
    return CandleCirq {
        Cs: cs,
        Cap: n,
        Fnt: 0,
        Cnt: 0,
    }
}

// This allows for over writing old values
func (cq *CandleCirq) Enqueue(c Candle) {
    cq.Cs[cq.Fnt] = c
    cq.Fnt = (cq.Fnt + 1) % cq.Cap
    if cq.Cnt < cq.Cap {
        cq.Cnt += 1
    }
}

func (cq *CandleCirq) Dequeue() (Candle, error) {
    var idx uint32
    
    if cq.Cnt == 0 {
        return Candle{}, fmt.Errorf("candle queue is empty")
    }
    
    idx = (cq.Cap - cq.Cnt + cq.Fnt) % cq.Cap // same as Fnt - Cnt mod Cap
    cq.Cnt -= 1

    return cq.Cs[idx], nil
}

func (cq CandleCirq) Copy() CandleCirq {
    var (
        cs      []Candle
        cnt     uint32
        idx     uint32
    )

    cs  = make([]Candle, cq.Cap)
    cnt = cq.Cnt
    idx = cq.Fnt

    for cnt > 0 {
        cs[idx] = cq.Cs[idx].Copy()
        idx = (idx + cq.Cap - 1) % cq.Cap
        cnt -= 1
    }

    return CandleCirq {
        Cs: cs,
        Fnt: cq.Fnt,
        Cnt: cq.Cnt,
        Cap: cq.Cap,
    }
}
