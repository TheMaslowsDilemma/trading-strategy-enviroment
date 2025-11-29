package handlers

import (
	"encoding/json"
	"tse-p3/ledger"
)


type SubscribeCmd struct {
	Name  string            `json:"name"`
	Etype ledger.EntityType `json:"entity_type"`
	Addr  string 			`json:"address"`
}

type UnsubscribeCmd struct {
	Etype	ledger.EntityType	`json:"entity_type"`
	Addr	string				`json:"address"`
}

type SearchCmd struct {
	Name string `json:"name"`
}

type SwapCmd struct {
	AmountIn   uint64 `json:"amount_in"`
	FromSymbol string `json:"from_symbol"`
	ToSymbol   string `json:"to_symbol"`
}

type CommandMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}