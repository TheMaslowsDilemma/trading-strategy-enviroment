package bots

import (
	"fmt"
	"tse-p3/traders"
	"tse-p3/strategy"
	"tse-p3/transactions"
)

type Bot struct {
	Name		string
	PendingTx	bool
	Id		uint64
	Strategy	strategies.Strategy
	Trader		*traders.Trader
}

func (bot *Bot) String() string {
	return fmt.Sprintf("{ name: %v; id: %v; pnd-tx: %v}", bot.Name, bot.Id, bot.PendingTx)
}

func (bot *Bot) NotificationHandler(res txs.TxResult) {
	bot.PendingTx = false
}
