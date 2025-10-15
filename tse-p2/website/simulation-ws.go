// simulation-ws.go
package website

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"tse-p2/candles"
	"tse-p2/simulation"
	"tse-p2/ledger"
	"tse-p2/exchange"
)

type hub struct {
	clients    map[*websocket.Conn]bool
	Broadcast chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
}

func newHub() *hub {
	return &hub{
		Broadcast: make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		clients:    make(map[*websocket.Conn]bool),
	}
}

func (h *hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
			}
		case message := <-h.Broadcast:
			for client := range h.clients {
				if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
					client.Close()
					delete(h.clients, client)
				}
			}
		}
	}
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func websocketHandler(w http.ResponseWriter, r *http.Request, sim *simulation.Simulation) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("WebSocket upgrade error: %v\n", err)
		return
	}

	Hub.register <- conn

	defer func() {
		Hub.unregister <- conn
		conn.Close()
	}()

	// Send initial candles on connection
	initialCandles, err := GetCandles(sim)
	if err == nil {
		data, _ := json.Marshal(map[string]interface{}{
			"type": "candles",
			"data": initialCandles,
		})
		conn.WriteMessage(websocket.TextMessage, data)
	}

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

			response := processCommand(msg.Command, sim)
			respData, _ := json.Marshal(map[string]string{
				"type":    "response",
				"message": response,
			})

			conn.WriteMessage(websocket.TextMessage, respData)

		} else if msg.Type == "init" {

			candles, err := GetCandles(sim)

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

func GetCandles(sim *simulation.Simulation) ([]candles.Candle, error) {
	sim.LedgerLock.Lock()
	defer sim.LedgerLock.Unlock()

	exItem, ok := sim.Ledger[sim.ExAddr]
	if !ok {
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