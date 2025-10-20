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
	var (
		t1  trader.Trader
		t2  trader.Trader
		t3  trader.Trader
	)

	s1 := strategy.SimpleStrategy {
		ShortInterval: 2,
		LongInterval: 6,
	}

	s2 := strategy.SimpleStrategy {
		ShortInterval: 7,
		LongInterval: 18,
	}

	s3 := strategy.SimpleStrategy {
		ShortInterval: 16,
		LongInterval: 42,
	}

	t1 = sim.createTrader(s1)
	t2 = sim.createTrader(s2)
	t3 = sim.createTrader(s3)

	go t1.Run(&sim.IsCanceled)
	go t2.Run(&sim.IsCanceled)
	go t3.Run(&sim.IsCanceled)
}


