package simulation

import (
	"tse-p3/strategy"
)

type BotDescriptor struct {
	Strat	strategies.Strategy
	Name	string
}

func (sim *Simulation) InitializeTraders() {
	var InterestingStrategies = []BotDescriptor {
		{
	        Strat: strategies.SimpleStrategy{ShortInterval: 3, LongInterval: 5},
			Name:  "simple-short",
	    },
	    {
       		Strat: strategies.SimpleStrategy{ShortInterval: 4, LongInterval: 6},
        	Name:  "simple-long",
    	},
    	{
			Strat: strategies.RandomStrategy{},
			Name: "rnd-001",
		},
		{
			Strat: strategies.RandomStrategy{},
			Name: "rnd-002",
		},
	}

	for _, is := range InterestingStrategies {
		botid := sim.AddBot(is.Name, is.Strat)
		bot := sim.Bots[botid]
		go bot.Run(&sim.CancelRequested, sim.GetCandles, sim.PlaceBotSwap, sim.GetWallet)
	}
}
