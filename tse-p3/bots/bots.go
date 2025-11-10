package bots

import (
	"fmt"
	"tse-p3/traders"
	"tse-p3/strategy"
	"tse-p3/transactions"
)

type Bot struct {
	Name		string
	Id			uint64
	Strategy	strategies.Strategy
	Trader		*traders.Trader
}

func (bot *Bot) String() string {
	return fmt.Sprintf("{ name: %v; id: %v }", bot.Name, bot.Id)
}

func (bot *Bot) NotificationHandler(res txs.TxResult) {
	fmt.Printf("[%v] tx %v\n", bot.Name, res)
}
