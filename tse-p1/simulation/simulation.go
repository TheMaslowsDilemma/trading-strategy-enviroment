package simulation

import (
	"fmt"
	"tse-p1/candles"
	"tse-p1/market"
	"tse-p1/strategy"
)


type Simulator struct {
	Balance float64
	Position float64
	Strategy strategy.Strategy
	TradeFee float64
	
}

func NewSimulator(initialBalance float64, strategy strategy.Strategy, tradeFee float64) *Simulator {
	return &Simulator{
		Balance: initialBalance,
		Position: 0.0,
		Strategy: strategy,
		TradeFee: tradeFee,
	}
}

func (sim *Simulator) Run(candles []candles.Candle) []float64 {
	var (
		action market.Action
		lastprice float64
		networth float64
		ns []float64
		i int
	)

	for i = 0; i < len(candles); i++ {
		action = sim.Strategy.Decide(candles, i)
		lastprice = candles[i].Close
		networth = sim.Position * lastprice + sim.Balance
		ns = append(ns, networth)

		if networth <= 0.001 {
			fmt.Printf("Bankrupcy at: %v\n", candles[i].Timestamp)
			return ns
		}

		switch action {
		case market.Buy:
			if sim.Balance > 0 {
				buy_usd := sim.Balance
				fee := buy_usd * sim.TradeFee
				rcv_amt := buy_usd - fee

				if (rcv_amt <= 0) {
					continue
				}

				sim.Balance  -= buy_usd
				sim.Position += rcv_amt / lastprice
			}
		case market.Sell:
			if sim.Position > 0 {
				proceeds := sim.Position * lastprice
				fee := proceeds * sim.TradeFee
				sim.Balance += proceeds - fee
				sim.Position = 0
			}
		}
	}
	return ns
}