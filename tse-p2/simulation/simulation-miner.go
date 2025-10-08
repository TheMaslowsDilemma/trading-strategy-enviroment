package simulation

// Standard Pattern for Simulation Entities:
// 1. infinite for loop
// 2. select to either end task & cleanup or iterate entity
// 3. Optional Delay
func (sim *Simulation) minerTask() {
	for {
		select {
			case <-sim.CancelChan:
				cleanMinerTask()
				return
			default:
				iterateMinerTask()
				break
		}
		// thread.Sleep
	}
}


func (sim *Simulation) iterateMinerTask() {

}

func (sim *Simulation) cleanMinerTask() {

} 
