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
	tkr = tokens.CreateTokenReserve(PrimarySymbol, 10146)

	fmt.Println(tkr)
}
