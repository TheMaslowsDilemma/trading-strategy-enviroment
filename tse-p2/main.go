package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "html/template"
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

var tmpl = template.Must(template.New("home").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Candle Audit Viewer</title>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/luxon"></script>
    <script src="https://cdn.jsdelivr.net/npm/chartjs-adapter-luxon"></script>
    <style>
        body { 
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; 
            background: linear-gradient(135deg, #1e3c72 0%, #2a5298 100%);
            color: white;
            margin: 0;
            padding: 20px;
            min-height: 100vh;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        h1 {
            text-align: center;
            margin-bottom: 30px;
            font-size: 2.5em;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.5);
        }
        #chart-container { 
            background: rgba(255,255,255,0.1);
            border-radius: 15px;
            padding: 20px;
            margin-bottom: 30px;
            backdrop-filter: blur(10px);
            box-shadow: 0 8px 32px rgba(0,0,0,0.3);
        }
        #command-input { 
            background: rgba(255,255,255,0.1);
            border-radius: 10px;
            padding: 20px;
            backdrop-filter: blur(10px);
            box-shadow: 0 8px 32px rgba(0,0,0,0.3);
        }
        #command {
            width: 60%;
            padding: 12px;
            border: none;
            border-radius: 5px;
            font-size: 16px;
            background: rgba(255,255,255,0.9);
            color: #333;
        }
        button {
            padding: 12px 24px;
            border: none;
            border-radius: 5px;
            background: #4CAF50;
            color: white;
            font-size: 16px;
            cursor: pointer;
            margin-left: 10px;
        }
        button:hover {
            background: #45a049;
        }
        #output {
            margin-top: 20px;
            background: rgba(0,0,0,0.3);
            border-radius: 10px;
            padding: 15px;
            max-height: 200px;
            overflow-y: auto;
            font-family: monospace;
            font-size: 14px;
            white-space: pre-wrap;
        }
        .stats {
            display: flex;
            justify-content: space-around;
            margin-bottom: 20px;
            background: rgba(255,255,255,0.1);
            padding: 15px;
            border-radius: 10px;
        }
        .stat-item {
            text-align: center;
        }
        .stat-value {
            font-size: 1.5em;
            font-weight: bold;
            color: #FFD700;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Trading Simulation</h1>
        
        <div class="stats">
            <div class="stat-item">
                <div class="stat-value" id="currentPrice">-</div>
                <div>Current Close</div>
            </div>
            <div class="stat-item">
                <div class="stat-value" id="candleCount">-</div>
                <div>Total Candles</div>
            </div>
            <div class="stat-item">
                <div class="stat-value" id="simDuration">-</div>
                <div>Simulation Time</div>
            </div>
        </div>

        <div id="chart-container">
            <canvas id="priceChart"></canvas>
        </div>
        
        <div id="command-input">
            <h3>Command Interface</h3>
            <input type="text" autocorrect="off" id="command" placeholder="Enter command (e.g., swap A B 0.5, getdur, help)" />
            <button onclick="sendCommand()">Execute</button>
            <div style="margin-top: 10px; font-size: 12px; opacity: 0.8;">
                Type 'help' for available commands
            </div>
        </div>
        
        <div id="output"></div>
    </div>

    <script>
        var ws = new WebSocket("ws://" + location.host + "/ws");
        var chart;
        var closePrices = [];

        ws.onopen = function() {
            console.log("WebSocket connected");
            ws.send(JSON.stringify({type: "init"}));
        };

        ws.onmessage = function(event) {
            var msg = JSON.parse(event.data);
            if (msg.type === "candles") {
                updateChart(msg.data);
                updateStats(msg.data);
            } else if (msg.type === "response") {
                document.getElementById("output").innerText += msg.message + "\n";
                document.getElementById("output").scrollTop = document.getElementById("output").scrollHeight;
            } else if (msg.type === "stats") {
                updateStats(msg.data);
            }
        };

        function sendCommand() {
            var cmd = document.getElementById("command").value;
            if (cmd.trim()) {
                ws.send(JSON.stringify({type: "command", command: cmd}));
                document.getElementById("command").value = "";
            }
        }

        // Allow Enter key to send command
        document.getElementById("command").addEventListener("keypress", function(e) {
            if (e.key === "Enter") {
                sendCommand();
            }
        });

        function updateChart(candles) {
            // Extract close prices and timestamps
            closePrices = candles.map(c => ({
                x: new Date(c.TimeStamp * 1000),
                y: c.Close
            })).reverse(); // Reverse to show chronological order

            if (!chart) {
                var ctx = document.getElementById("priceChart").getContext("2d");
                chart = new Chart(ctx, {
                    type: "line",
                    data: {
                        datasets: [{
                            label: "Close Price",
                            data: closePrices,
                            borderColor: "#00D4FF",
                            backgroundColor: "rgba(0, 212, 255, 0.1)",
                            borderWidth: 3,
                            fill: true,
                            tension: 0.1,
                            pointRadius: 0,
                            pointHoverRadius: 6,
                            pointBackgroundColor: "#00D4FF",
                            pointHoverBackgroundColor: "#FFD700"
                        }]
                    },
                    options: {
                        responsive: true,
                        maintainAspectRatio: false,
                        interaction: {
                            intersect: false,
                            mode: 'index'
                        },
                        plugins: {
                            legend: {
                                labels: {
                                    font: { size: 14 },
                                    color: 'white'
                                }
                            },
                            tooltip: {
                                backgroundColor: 'rgba(0,0,0,0.8)',
                                titleColor: 'white',
                                bodyColor: 'white',
                                borderColor: '#00D4FF',
                                borderWidth: 1,
                                callbacks: {
                                    label: function(context) {
                                        return 'Price: $' + context.parsed.y.toFixed(4);
                                    }
                                }
                            }
                        },
                        scales: {
                            x: {
                                type: "time",
                                time: {
                                    unit: "second",
                                    displayFormats: {
                                        second: 'HH:mm:ss'
                                    }
                                },
                                grid: {
                                    color: "rgba(255,255,255,0.1)"
                                },
                                ticks: {
                                    color: "white",
                                    maxTicksLimit: 10
                                },
                                title: {
                                    display: true,
                                    text: "Time",
                                    color: "white"
                                }
                            },
                            y: {
                                grid: {
                                    color: "rgba(255,255,255,0.1)"
                                },
                                ticks: {
                                    color: "white",
                                    callback: function(value) {
                                        return '$' + value.toFixed(4);
                                    }
                                },
                                title: {
                                    display: true,
                                    text: "Price (USD)",
                                    color: "white"
                                }
                            }
                        },
                        elements: {
                            point: {
                                hoverBorderWidth: 3
                            }
                        },
                        animation: {
                            duration: 500,
                            easing: 'easeInOutQuart'
                        }
                    }
                });
                document.getElementById("chart-container").style.height = "500px";
            } else {
                chart.data.datasets[0].data = closePrices;
                chart.update('none'); // Fast update without animation for real-time
            }
        }

        function updateStats(candles) {
            if (candles.length > 0) {
                var latest = candles[candles.length - 1];
                document.getElementById("currentPrice").textContent = '$' + latest.Close.toFixed(4);
                document.getElementById("candleCount").textContent = candles.length;
            }
        }
    </script>
</body>
</html>
`))

// wsHandler handles WebSocket connections
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