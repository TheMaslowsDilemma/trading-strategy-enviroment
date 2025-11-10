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

	// Add Entities //
//	(&s).AddUser("tom", 1234)
	(&s).InitializeTraders()

	

	for _, bot := range s.Bots {
		fmt.Printf("[%v] networth: %v\n", bot.Name, bot.Trader.GetNetworth(s.GetPrice, s.GetWallet))
	}

	time.Sleep(60 * time.Second)
	s.CancelRequested = true
	
	fmt.Println(s.String())

	for _, bot := range s.Bots {
		fmt.Printf("[%v] networth: %v\n", bot.Name, bot.Trader.GetNetworth(s.GetPrice, s.GetWallet))
	}
}
