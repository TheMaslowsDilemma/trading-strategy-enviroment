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
	(&s).AddUser("tom", 1234)
	(&s).InitializeTraders()

	(&s).PlaceUserSwap(1234, "usd", "tse", 10)
	(&s).PlaceUserSwap(1234, "usd", "tse", 10)
	(&s).PlaceUserSwap(1234, "usd", "tse", 10)
	(&s).PlaceUserSwap(1234, "usd", "tse", 10)
	time.Sleep(1 * time.Second)
	(&s).PlaceUserSwap(1234, "usd", "tse", 10)
	(&s).PlaceUserSwap(1234, "usd", "tse", 11)
	time.Sleep(1 * time.Second)
	(&s).PlaceUserSwap(1234, "tse", "usd", 1000)
	(&s).PlaceUserSwap(1234, "tse", "usd", 30)
	(&s).PlaceUserSwap(1234, "tse", "usd", 1000)
	(&s).PlaceUserSwap(1234, "tse", "usd", 30)
	time.Sleep(1 * time.Second)
	(&s).PlaceUserSwap(1234, "usd", "tse", 10)
	(&s).PlaceUserSwap(1234, "usd", "tse", 11)
	(&s).PlaceUserSwap(1234, "tse", "usd", 100)
	(&s).PlaceUserSwap(1234, "tse", "usd", 30)
	time.Sleep(1 * time.Second)
	(&s).PlaceUserSwap(1234, "usd", "tse", 10)
	(&s).PlaceUserSwap(1234, "usd", "tse", 10)
	(&s).PlaceUserSwap(1234, "usd", "tse", 10)
	(&s).PlaceUserSwap(1234, "usd", "tse", 10)
	(&s).PlaceUserSwap(1234, "usd", "tse", 10)
	time.Sleep(1 * time.Second)
	(&s).PlaceUserSwap(1234, "usd", "tse", 11)
	(&s).PlaceUserSwap(1234, "tse", "usd", 1000)
	(&s).PlaceUserSwap(1234, "tse", "usd", 30)
	(&s).PlaceUserSwap(1234, "tse", "usd", 1000)
	(&s).PlaceUserSwap(1234, "tse", "usd", 30)
	time.Sleep(1 * time.Second)
	(&s).PlaceUserSwap(1234, "usd", "tse", 10)
	(&s).PlaceUserSwap(1234, "usd", "tse", 11)
	(&s).PlaceUserSwap(1234, "tse", "usd", 100)
	time.Sleep(1 * time.Second)
	(&s).PlaceUserSwap(1234, "tse", "usd", 30)
	(&s).PlaceUserSwap(1234, "usd", "tse", 10)
	time.Sleep(1 * time.Second)
	(&s).PlaceUserSwap(1234, "usd", "tse", 10)
	(&s).PlaceUserSwap(1234, "usd", "tse", 10)
	time.Sleep(1 * time.Second)
	(&s).PlaceUserSwap(1234, "usd", "tse", 10)
	(&s).PlaceUserSwap(1234, "usd", "tse", 10)
	time.Sleep(1 * time.Second)
	(&s).PlaceUserSwap(1234, "usd", "tse", 11)
	(&s).PlaceUserSwap(1234, "tse", "usd", 1000)
	(&s).PlaceUserSwap(1234, "tse", "usd", 30)
	(&s).PlaceUserSwap(1234, "tse", "usd", 1000)
	(&s).PlaceUserSwap(1234, "tse", "usd", 30)
	time.Sleep(1 * time.Second)
	(&s).PlaceUserSwap(1234, "usd", "tse", 10)
	(&s).PlaceUserSwap(1234, "usd", "tse", 11)
	(&s).PlaceUserSwap(1234, "tse", "usd", 100)
	time.Sleep(1 * time.Second)

	(&s).PlaceUserSwap(1234, "tse", "usd", 30)
	time.Sleep(60 * time.Second)
	s.CancelRequested = true
	fmt.Println(s.String())
}
