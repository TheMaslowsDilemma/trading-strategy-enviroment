package trader

import (
	"time"
    "tse-p2/candles"
    "tse-p2/ledger"
)

const traderIterationPeriod = 500 * time.Millisecond

func (t *Trader) Run(isCanceled *bool) {
	var tick uint64 = 0

    for {
    	if *isCanceled {
            return
        }

        t.iterate(tick)

        tick += 1
        time.Sleep(traderIterationPeriod)
    }
}

func (t *Trader) iterate(tick uint64) {
    var (
        cs  []candles.Candle
        tx  ledger.Tx
        err error
    )

    cs, err = t.candleFetcher()
    if err != nil {
        // NOTE we could log some error here...
        return
    }
    tx, _ = t.createTransaction(cs)
    if tx == nil {
        // NOTE we could log debug here saying we have no tx to make
        return
    }
    t.txPlacer(tx)
}
