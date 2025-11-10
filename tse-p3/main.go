package main

import (
	"fmt"
	"time"
	"tse-p3/simulation"
)

func main() {
	fmt.Println("--- Trading Stategy Environment: Part Three ---")
	
	s := simulation.NewSimulation()
	(&s).InitializeTraders()
	go (&s).Run()

	

	for _, bot := range s.Bots {
		fmt.Printf("[%v] networth: %v\n", bot.Name, bot.Trader.GetNetworth(s.GetPrice, s.GetWallet))
	}

	time.Sleep(10 * time.Second)
	s.CancelRequested = true
	
	fmt.Println(s.String())

	for _, bot := range s.Bots {
		fmt.Printf("[%v] networth: %v\n", bot.Name, bot.Trader.GetNetworth(s.GetPrice, s.GetWallet))
	}
}
