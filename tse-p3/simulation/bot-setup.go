package simulation

import (
	"fmt"
	"math/rand"
	"tse-p3/strategy"
)

type BotDescriptor struct {
	Strat strategies.Strategy
	Name  string
}

// add random bot to get


func (sim *Simulation) InitializeTraders() {
	randbot_d := BotDescriptor {
		Name: "random-01",
		Strat: strategies.RandomStrategy {},
	}

	randbotid := sim.AddBot(randbot_d.Name, randbot_d.Strat)
	randbot := sim.Bots[randbotid]
	go randbot.Run(&sim.CancelRequested, sim.GetCandles, sim.PlaceBotSwap, sim.GetWallet)
	
	randbot_d = BotDescriptor {
		Name: "random-02",
		Strat: strategies.RandomStrategy {},
	}

	randbotid = sim.AddBot(randbot_d.Name, randbot_d.Strat)
	randbot = sim.Bots[randbotid]
	go randbot.Run(&sim.CancelRequested, sim.GetCandles, sim.PlaceBotSwap, sim.GetWallet)
	

	for i := 0; i < 200; i++ {
		short := 5 + rand.Intn(46)
		long := short + 15 + rand.Intn(150)

		bd := BotDescriptor{
			Strat: strategies.SimpleStrategy{
				ShortInterval: short,
				LongInterval:  long,
			},
			Name: fmt.Sprintf("simple-%04d", i), // e.g., simple-0001, simple-0399
		}

		botid := sim.AddBot(bd.Name, bd.Strat)
		bot := sim.Bots[botid]
		go bot.Run(&sim.CancelRequested, sim.GetCandles, sim.PlaceBotSwap, sim.GetWallet)
	}
}