package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "math/rand"
    "net/http"
    "strconv"
    "strings"
    "sync"
    "time"

    "github.com/gorilla/websocket"
    "tse-p2/ledger"
    "tse-p2/simulation"
    "tse-p2/candles"
    "tse-p2/exchange"
)

var (
    upgrader = websocket.Upgrader{
        CheckOrigin: func(r *http.Request) bool {
            return true // Allow all origins for simplicity; adjust for production
        },
    }
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
    rand.Seed(rsd)
    fmt.Println("---------------------------------------")
    fmt.Println("Trading Strategy Environment: Part Two")
    fmt.Printf("seed: %v\n", rsd)
    fmt.Println("---------------------------------------")



    sim, err = simulation.CreateSimulation()

    fmt.Printf("init:\n--> user-wallet: %v\n--> mainexchange: %v\n\n", sim.CliWallet, sim.ExAddr)

    if err != nil {
        fmt.Println(err)
        return
    }

    go sim.Run()

    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        wsHandler(w, r, sim)
    })
    fmt.Printf("Starting web server on %s\n", addr)
    if err := http.ListenAndServe(addr, nil); err != nil {
        fmt.Printf("Web server error: %v\n", err)
    }

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    tmpl.Execute(w, nil)
}

func wsHandler(w http.ResponseWriter, r *http.Request, sim *simulation.Simulation) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Printf("WebSocket upgrade error: %v\n", err)
        return
    }
    defer conn.Close()

    // Mutex for safe access if multiple clients, but for now simple
    var mu sync.Mutex

    // Periodically send candle updates
    ticker := time.NewTicker(1 * time.Second) // Update every second; adjust as needed
    defer ticker.Stop()

    go func() {
        for range ticker.C {
            mu.Lock()
            candles, err := getCandles(sim)
            mu.Unlock()
            if err != nil {
                fmt.Printf("Error getting candles: %v\n", err)
                continue
            }
            data, _ := json.Marshal(map[string]interface{}{
                "type": "candles",
                "data": candles,
            })
            conn.WriteMessage(websocket.TextMessage, data)
        }
    }()

    for {
        _, message, err := conn.ReadMessage()
        if err != nil {
            fmt.Printf("WebSocket read error: %v\n", err)
            break
        }

        var msg struct {
            Type    string `json:"type"`
            Command string `json:"command"`
        }
        if err := json.Unmarshal(message, &msg); err != nil {
            continue
        }

        if msg.Type == "command" {
            // Process command similar to RunCLI
            response := processCommand(msg.Command, sim)
            respData, _ := json.Marshal(map[string]string{
                "type":    "response",
                "message": response,
            })
            conn.WriteMessage(websocket.TextMessage, respData)
        } else if msg.Type == "init" {
            // Send initial candles
            mu.Lock()
            candles, err := getCandles(sim)
            mu.Unlock()
            if err == nil {
                data, _ := json.Marshal(map[string]interface{}{
                    "type": "candles",
                    "data": candles,
                })
                conn.WriteMessage(websocket.TextMessage, data)
            }
        }
    }
}

// getCandles retrieves the candle history and current candle
func getCandles(sim *simulation.Simulation) ([]candles.Candle, error) {
    sim.LedgerLock.Lock()
    defer sim.LedgerLock.Unlock()

    exItem, ok := sim.Ledger[sim.ExAddr]
    if !ok  {
        return nil, fmt.Errorf("exchange not found on ledger")
    }
    ex, ok := exItem.(exchange.ConstantProductExchange)
    if !ok {
        return nil, fmt.Errorf("invalid exchange type")
    }

    auditItem, ok := sim.Ledger[ex.CndlAddr]
    if !ok {
        return nil, fmt.Errorf("candle audit not found on ledger")
    }
    audit, ok := auditItem.(candles.CandleAudit)
    if !ok {
        return nil, fmt.Errorf("invalid candle audit type")
    }

    candles := audit.CandleHistory.CandlesInOrder()
    candles = append(candles, audit.CurrentCandle)

    return candles, nil
}

func processCommand(s string, sim *simulation.Simulation) string {
    s = strings.Trim(s, " \n")
    if s == "q" {
        sim.IsCanceled = true
        return "signaling sim shutdown"
    } else if s == "getdur" {
        return fmt.Sprintf("sim duration: %v", sim.RunningDur)
    } else if s == "help" {
        return HelpString
    } else if strings.HasPrefix(s, "getitem ") {
        idstr := strings.TrimSpace(s[len("getitem"):])
        id, err := strconv.ParseUint(idstr, 10, 64)
        if err != nil {
            return fmt.Sprintf("invalid id for getitem: %v", err)
        }
        listr, err := sim.GetLedgerItemString(ledger.LedgerAddr(id))
        if err != nil {
            return fmt.Sprintf("failed to get item from ledger: %v", err)
        }
        return listr
    } else if strings.HasPrefix(s, "swap ") {
        parts := strings.Fields(s[len("swap "):])
        if len(parts) != 3 {
            return "swap requires 3 arguments: from-symbol to-symbol confidence"
        }
        from := parts[0]
        to := parts[1]
        cnfd, err := strconv.ParseFloat(parts[2], 64)
        if err != nil {
            return fmt.Sprintf("invalid confidence for swap: %v", err)
        }
        err = sim.PlaceUserTrade(from, to, cnfd)
        if err != nil {
            return err.Error()
        }
        return fmt.Sprintf("swap order placed: %v %s to %s", cnfd, from, to)
    } else {
        return fmt.Sprintf("unrecognized command \"%s\"", s)
    }
}

var HelpString string = "\n\t\"help\": list of commands" +
    "\n\t\"getdur\": get time duration since simulation start" +
    "\n\t\"getitem <id>\": print value of ledger for <id>" +
    "\n\t\"swap\" <from-symbol> <to-symbol> <confidence [0,1]>"
