package simulation

import (
	"tse-p2/trader"
	"tse-p2/strategy"
	"tse-p2/wallet"
)

func (sim *Simulation) createTrader(s strategy.Strategy) trader.Trader {
	walletAddr := wallet.InitDefaultWallet(&sim.Ledger)

	return trader.CreateTrader(
		s,
        10,
        walletAddr,
        sim.ExAddr,
        "usd",
        "eth",
        sim.GetCandles,
    	sim.ledgerLookup,
    	sim.placeTx,
    )
}

func (sim *Simulation) initializeTraders() {
	var InterestingStrategies = []struct {
	    Strat strategy.Strategy
	    Name  string
	}{
	    {
	        Strat: strategy.SimpleStrategy{ShortInterval: 5, LongInterval: 20},
	        Name:  "SimpleStrategy-AggressiveShort",
	    },
	    {
	        Strat: strategy.MomentumStrategy{Lookback: 14, Threshold: 0.005},
	        Name:  "MomentumStrategy-Sensitive",
	    },
	    {
	        Strat: strategy.VolatilityBreakoutStrategy{ATRPeriod: 10, Multiplier: 2.0},
	        Name:  "VolatilityBreakoutStrategy-Frequent",
	    },
	    {
	        Strat: strategy.MeanReversionStrategy{SMAPeriod: 50, Deviation: 0.02},
	        Name:  "MeanReversionStrategy-Stable",
	    },
	    {
	        Strat: strategy.RandomWalkStrategy{BuyProb: 0.4, SellProb: 0.2},
	        Name:  "RandomWalkStrategy-BullBiased",
	    },
	    // Iter 2
		{
	        Strat: strategy.SimpleStrategy{ShortInterval: 25, LongInterval: 60},
	        Name:  "SimpleStrategy-AggressiveShort",
	    },
	    {
	        Strat: strategy.MomentumStrategy{Lookback: 50, Threshold: 0.009},
	        Name:  "MomentumStrategy-Sensitive",
	    },
	    {
	        Strat: strategy.VolatilityBreakoutStrategy{ATRPeriod: 30, Multiplier: 0.4},
	        Name:  "VolatilityBreakoutStrategy-Frequent",
	    },
	    {
	        Strat: strategy.MeanReversionStrategy{SMAPeriod: 100, Deviation: 0.01},
	        Name:  "MeanReversionStrategy-Stable",
	    },
	    {
	        Strat: strategy.RandomWalkStrategy{BuyProb: 0.4, SellProb: 0.2},
	        Name:  "RandomWalkStrategy-BullBiased",
	    },
	}

	for _, is := range InterestingStrategies {
		t := sim.createTrader(is.Strat)
		go t.Run(&sim.IsCanceled)
	}
}