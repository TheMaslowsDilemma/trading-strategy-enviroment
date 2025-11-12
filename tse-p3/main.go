package main

import (
	"fmt"
	"time"
	"tse-p3/simulation"
	"tse-p3/globals"
)



func main() {
	fmt.Println("--- Trading Stategy Environment: Part Three ---")
	
	s := simulation.NewSimulation()
	(&s).InitializeTraders()
	go (&s).Run()

	
	for _, bot := range s.Bots {
		fmt.Printf("[%v] networth: %v\n", bot.Name, bot.Trader.GetNetworth(s.GetPrice, s.GetWallet))
	}

	time.Sleep(15 * time.Second)
	s.CancelRequested = true
	
	fmt.Println(s.String())

	fmt.Println("\n--- Simulation Results ---")
	
	nws := make([]float64, 0)
	var total_nws float64 = 0.0

	for _, bot := range s.Bots {
		var nw = bot.Trader.GetNetworth(s.GetPrice, s.GetWallet)
		total_nws += nw
		fmt.Printf("\t[%v] : %v : %v \n", bot.Name, bot.Trader.String(s.GetWallet), nw)
		nws = append(nws, nw)
	}

	fmt.Printf("total: %v\n", total_nws)
	fmt.Printf("\t[exchange] : %v\n", s.GetExchange(globals.USDSymbol, globals.TSESymbol))
}
