package simulation

import (
	"tse-p3/strategy"
)

func (sim *Simulation) InitializeTraders() {
	var InterestingStrategies = []struct {
	    Strat strategies.Strategy
	    Name  string
	}{
	    {
	        Strat: strategies.SimpleStrategy{ShortInterval: 3, LongInterval: 5},
	        Name:  "simple-short",
	    },
		{
	        Strat: strategies.SimpleStrategy{ShortInterval: 15, LongInterval: 20},
	        Name:  "simple-long",
	    },
	}

	for _, is := range InterestingStrategies {
		botid := sim.AddBot(is.Name, is.Strat)
		bot := sim.Bots[botid]
		go bot.Run(&sim.CancelRequested, sim.GetCandles, sim.PlaceBotSwap, sim.getWallet)
	}
}