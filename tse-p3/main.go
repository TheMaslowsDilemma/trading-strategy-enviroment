package main

import (
	"fmt"
	"tse-p3/simulation"
)

func main() {
	fmt.Println("--- Trading Stategy Environment: Part Three ---")
	
	s := simulation.NewSimulation()

	fmt.Println(s.String())
}
