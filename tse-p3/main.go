package main

import (
	"fmt"
	"tse-p3/tokens"
)

func main() {
	fmt.Println("--- Trading Stategy Environment: Part Three ---")
	
	var (
		PrimarySymbol 	string
		tkr		tokens.TokenReserve
	)

	PrimarySymbol = "usd"
	tkr = tokens.CreateTokenReserve(101245, PrimarySymbol)

	fmt.Println(tkr)
}
