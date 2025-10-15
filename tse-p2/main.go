package main

import (
    "encoding/json"
    "flag"
    "os"
    "fmt"
    "math/rand"
    "strconv"
    "time"

    "tse-p2/simulation"
    "tse-p2/website"
)


func main() {
    var (
        rsd     int64
        sim     *simulation.Simulation
        err     error
        addr    string
    )

    flag.StringVar(&addr, "addr", ":8080", "http service address")
    flag.Parse()

    rsd = time.Now().UnixNano()

    // -- optionally supply the random seed in the command line -- //
    if len(flag.Args()) > 0 {
        if seed, err := strconv.ParseInt(flag.Args()[0], 10, 64); err == nil {
            rsd = seed
        } else {
            fmt.Fprintf(os.Stderr, "Invalid random seed arg: %v, continuing with unix timestamp based\n", err)
        }
    }

    rand.Seed(rsd)
    sim, err = simulation.CreateSimulation()
    if err != nil {
        fmt.Println(err)
        return
    }

    website.Initialize(addr, sim)

    sim.CandleNotifier = func() {
        candles, err := website.GetCandles(sim)
        if err != nil {
            return
        }
        data, err := json.Marshal(map[string]interface{}{
            "type": "candles",
            "data": candles,
        })
        if err != nil {
            return
        }
        website.Hub.Broadcast <- data
    }

    go sim.Run()

    fmt.Println("---------------------------------------")
    fmt.Println("Trading Strategy Environment: Part Two")
    fmt.Printf("- random-seed   : %v\n", rsd)
    fmt.Printf("- wallet-addr   : %v\n", sim. CliWallet)
    fmt.Printf("- exchange-addr : %v\n", sim.ExAddr)
    fmt.Printf("- website at    : %v\n", website.Address)
    fmt.Println("---------------------------------------")
    website.Begin()
}