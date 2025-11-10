package simulation

import (
	"time"
	"tse-p3/miner"
	"tse-p3/globals"
)

func (sim *Simulation) Run() {
	var start, tick	uint64

	start = uint64(time.Now().Unix())
	for {
		if sim.CancelRequested {
			return
		}

		tick = uint64(time.Now().Unix()) - start
		sim.MinerTask(tick)
		time.Sleep(globals.TimeBetweenBlocks)
	}
}

func (sim *Simulation) MinerTask(tick uint64) {
	var (
		err		error
	)

	_, err = miner.NextBlock(tick, &sim.MemoryPool, &sim.ScndLedger)
	if err != nil {
		// TODO push err log to sim
		return
	}

	// NOTE this is where we bring delta into the main ledger
	sim.LedgerLock.Lock()
	(&sim.MainLedger).Merge(sim.ScndLedger)
	sim.LedgerLock.Unlock()
}
