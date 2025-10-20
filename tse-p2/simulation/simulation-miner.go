package simulation

import (
    "time"
    "tse-p2/ledger"
)

const timeBetweenBlocks = 100 * time.Millisecond

func (sim *Simulation) minerTask() {
    var start uint64 = uint64(time.Now().Unix())
    for {
    	if sim.IsCanceled {
            return
        }
        sim.iterateMinerTask( (uint64(time.Now().Unix()) - start))
        time.Sleep(timeBetweenBlocks)
    }
}

func (sim *Simulation) iterateMinerTask(tick uint64) {
    var (
        ftcount uint
        err     error
    )

    
    ftcount, err = sim.MainMiner.MineNextBlock(tick, &sim.MemoryPool)
    if err != nil {
        // TODO push err log to sim
    }

    
    ledger.Merge(&sim.Ledger, sim.MainMiner.BackLedger)
    
    if sim.CandleNotifier != nil && ftcount != 0 {
        sim.CandleNotifier()
    }
}
