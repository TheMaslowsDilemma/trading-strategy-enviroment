package main

import (
	"fmt"
	"time"
	"tse-p3/simulation"
)

func main() {
	fmt.Println("--- Trading Stategy Environment: Part Three ---")
	
	s := simulation.NewSimulation()
	
	go (&s).Run()

	time.Sleep(5 * time.Second)
	(&s).AddUser()
	s.CancelRequested = true
	fmt.Println(s.String())
}
