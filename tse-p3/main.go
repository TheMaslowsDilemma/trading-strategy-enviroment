package main

import (
	"fmt"
	"os"
	"tse-p3/simulation"
	"tse-p3/website"
	"tse-p3/db"
)



func main() {

	fmt.Println("--- Trading Stategy Environment: Part Three ---")
	db.Init()
	s := simulation.NewSimulation()
    website.Initialize(&s)
	(&s).InitializeTraders()
	go (&s).Run()

	port := os.Getenv("PORT")
	if port == "" {
	    port = "8080"
	}
	website.Begin(":8080")
}
