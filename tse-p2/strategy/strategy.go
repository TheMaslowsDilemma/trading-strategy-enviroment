package strategy

import (
    "tse-p2/candles"
)

type Action uint8
const (
    Buy = iota
    Sell
    Hold
)


type Strategy interface {
   Decide(cs []candles.Candle) (Action, float64) // cs -> Action, Confidence [0,1]
}
