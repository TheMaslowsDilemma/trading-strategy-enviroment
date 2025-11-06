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
	(&s).AddUser("tom", 1234)
	(&s).PlaceUserSwap(1234, "usd", "tse", 10)
	(&s).PlaceUserSwap(1234, "usd", "tse", 10)
	(&s).PlaceUserSwap(1234, "usd", "tse", 10)
	(&s).PlaceUserSwap(1234, "usd", "tse", 10)
	(&s).PlaceUserSwap(1234, "usd", "tse", 10)
	(&s).PlaceUserSwap(1234, "usd", "tse", 10)
	(&s).PlaceUserSwap(1234, "usd", "tse", 10)
	(&s).PlaceUserSwap(1234, "tse", "usd", 10)
	(&s).PlaceUserSwap(1234, "tse", "usd", 10)
	(&s).PlaceUserSwap(1234, "tse", "usd", 10)
	(&s).PlaceUserSwap(1234, "tse", "usd", 10)
	(&s).PlaceUserSwap(1234, "tse", "usd", 10)
	(&s).PlaceUserSwap(1234, "tse", "usd", 10)
	time.Sleep(5 * time.Second)
	s.CancelRequested = true
	fmt.Println(s.String())
}
