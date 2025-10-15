package simulation

import (
    "time"
    "tse-p2/ledger"
)

const timeBetweenMinerIterations = 100 * time.Millisecond
const timeBetweenBlocks = 500 * time.Millisecond

func (sim *Simulation) minerTask() {
    var tick uint64 = 0

    for {
    	if sim.IsCanceled {
            return
        }
        sim.iterateMinerTask(tick)
        tick += 1
        time.Sleep(timeBetweenMinerIterations)
    }
}

func (sim *Simulation) iterateMinerTask(tick uint64) {
    var (
        ftcount uint
        err     error
    )

    // TODO ? add CLI interations here
    if time.Since(sim.MainMiner.LastBlockTime) >= timeBetweenBlocks {
        ftcount, err = sim.MainMiner.MineNextBlock(tick, &sim.MemoryPool)
        if err != nil {
            // TODO push err log to sim
        }

        if sim.CandleNotifier != nil && ftcount != 0 {
            sim.CandleNotifier()
        }

        ledger.Merge(&sim.Ledger, sim.MainMiner.BackLedger)
    }
}
