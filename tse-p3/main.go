package main

import (
	"fmt"
	"time"
	"tse-p3/simulation"
	"tse-p3/website"
	"tse-p3/globals"
)



func main() {

	fmt.Println("--- Trading Stategy Environment: Part Three ---")
	
	s := simulation.NewSimulation()
    website.Initialize(&s)
	(&s).InitializeTraders()
	go (&s).Run()
	website.Begin(":8080")
}
