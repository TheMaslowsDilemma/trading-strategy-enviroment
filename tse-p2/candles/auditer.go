package candles

/***
TODO this is responsible for tracking prices
and generating candles. it watches the ledger
after each block is applied -- so it technically
will be used within the minerTask. it should
keep track of current Candle, which is under
construction until it is commited to history
***/
type CandleAuditer struct {
    CurrentCandle       Candle
    CandleHistory       []Candle
}
