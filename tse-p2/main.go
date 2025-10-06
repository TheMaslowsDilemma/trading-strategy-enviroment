package main

import (
    "os"
    "fmt"
    "time"
    "bufio"
    "strings"
    "strconv"
    "tse-p2/ledger"
    "tse-p2/wallet"
    "tse-p2/simulation"
)

func main() {
    var (
        rdr *bufio.Reader
        rch chan int
        sim *simulation.Simulation
        dur int64
        err error
    )

    fmt.Println("Trading Strategy Environment: Part Two")

    if len(os.Args) != 2 {
        fmt.Println("usage: go run . <duration-in-seconds>")
        return
    }
    
    dur, err = strconv.ParseInt(os.Args[1], 10, 64)
    if err != nil {
        fmt.Printf("Failed to parse duration: %v\n", err)
    }

    sim, err = simulation.CreateSimulation(time.Duration(dur) * time.Second)
    wallet1 := wallet.Wallet{ TraderId: 69, Reserves: make([]ledger.LedgerAddr, 2) }
    sim.AddLedgerItem(4, wallet1)

    if err != nil {
        fmt.Println(err)
        return
    }
    
    rch = make(chan int, 1)
    rch <- 0
    rdr = bufio.NewReader(os.Stdin)

    fmt.Printf("Simulation Created with duration %v\n", sim.Dur)
    fmt.Println("Simulation Starting")
    
    go sim.Run()
    for {
        select {
            case <-sim.CancelChan:
                sim.CancelChan <- 0
                fmt.Println("Simulation Complete")
                return
            default:
                RunCLI(&rch, rdr, sim)
                break
        }
    }
}

func RunCLI(rch *(chan int), rdr *bufio.Reader, sim *simulation.Simulation) {
    var c int
    select {
        case c = <-(*rch):
            RunUserCLI(c, rdr, sim)
            *rch <- c + 1
        default:
            return
    }
}

func RunUserCLI(c int, rdr *bufio.Reader, sim *simulation.Simulation) {
    var (
        sc      chan byte
        s       string
        e       error
    )

    fmt.Printf("(%v) > ", c)

    sc = make(chan byte, 1)
    go func() {
        s, e = rdr.ReadString('\n')
        sc <- 0
    }()

    select {
        case <-sc:
            break
        case <-sim.CancelChan:
            fmt.Println()
            return
    }
    if e != nil {
        fmt.Printf("(%v) >> Bad Read\n", c)
    }

    s = strings.Trim(s, " \n")
    if s == "q" {
        fmt.Printf("(%v) >> signaling sim shutdown\n", c)
        sim.CancelChan <- 1
    } else if s == "getdur" {
        fmt.Printf("(%v) >> sim duration: %v\n", c, sim.CurrentDur)
    } else if s == "help" {
        fmt.Printf("(%v) >> %s\n", c, HelpString)
    } else if strings.HasPrefix(s, "getitem") {
        var (
            idstr   string
            id      uint64
            listr   string
        )

        idstr = strings.TrimSpace(s[len("getitem"):])
        id, e = strconv.ParseUint(idstr, 10, 64)
        if e != nil {
            fmt.Printf("(%v) >> invalid id for getitem: %v\n", c, e)
            return
        }
        listr, e = sim.GetLedgerItemString(ledger.LedgerAddr(id))
        if e != nil {
            fmt.Printf("(%v) >> failed to get item from ledger: %v\n", c, e)
                return
        }
        fmt.Printf("(%v) >> %s\n", c, listr)
    } else {
        fmt.Printf("(%v) >> unrecognized command \"%s\"\n", c, s)
    }
}

var HelpString string = "\n\t\"help\": list of commands" +
    "\n\t\"getdur\": get time duration since simulation start" +
    "\n\t\"getitem <id>\": print value of ledger for <id>"
