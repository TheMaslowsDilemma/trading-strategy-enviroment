package simulation

import (
    "time"
    "tse-p2/ledger"
)

const timeBetweenMinerIterations = 50 * time.Millisecond
const timeBetweenBlocks = 250 * time.Millisecond

func (sim *Simulation) minerTask() {
    for {
	select {
	    case <-sim.CancelChan:
                return
	    default:
	        sim.iterateMinerTask()
	        break
	}
        time.Sleep(timeBetweenMinerIterations)
    }
}

func (sim *Simulation) iterateMinerTask() {
    // TODO add responses to external commands
    // this would happen before doing any ledger update
    var err error

    if time.Since(sim.MainMiner.LastBlockTime) >= timeBetweenBlocks {
        err = sim.MainMiner.MineNextBlock(&sim.MemoryPool)
        if err != nil {
            // TODO push err log to sim            
        }

        ledger.Merge(sim.Ledger, sim.MainMiner.BackLedger)
    }
}
