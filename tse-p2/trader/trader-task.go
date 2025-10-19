package trader

import (
	"time"
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
    
}