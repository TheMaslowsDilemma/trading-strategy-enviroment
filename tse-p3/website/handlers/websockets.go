package handlers

import (
	"fmt"
	"strconv"
	"context"
	"net/http"
	"encoding/json"

	"tse-p3/ledger"
	"tse-p3/users"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // tighten in production!
		},
	}
)

// Command structures
type SubscribeCmd struct {
	Name  string            `json:"name"`
	Etype ledger.EntityType `json:"entity_type"`
	Addr  string 			`json:"address"`
}

type UnsubscribeCmd struct {
	Etype ledger.EntityType `json:"entity_type"`
	Addr  ledger.Addr       `json:"address"`
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
	Type string          `json:"type"` // "subscribe" | "unsubscribe" | "search" | "swap"
	Data json.RawMessage `json:"data"`
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	var (
		conn		*websocket.Conn
		ctx			context.Context
		user		users.User
		userOK		bool
		msgType		int
		message		[]byte
		cmd			CommandMessage
		value		interface{}
		err			error
	)

	// 1. Upgrade connection
	conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("WebSocket upgrade failed: %v\n", err)
		return
	}
	defer conn.Close()

	// 2. Extract user id from context
	ctx = r.Context()
	value = ctx.Value("user")
	if value != nil {
		user, userOK = value.(users.User)
		if !userOK {
			fmt.Println("user value could not parse")
			conn.WriteJSON(map[string]string{"error": "unauthenticated"})
			return
		}
	} else {
		conn.WriteJSON(map[string]string{"error": "no user found"})
		return
	}

	fmt.Println("websocket added subscriptions")

	// -- initialize user data
	_ = conn.WriteJSON(map[string]interface{}{
		"type": "initialize",
		"data": user,
	})

	fmt.Println("websocket wrote welcome")


	// 4. Main message loop
	for {
		msgType, message, err = conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("WebSocket error: %v\n", err)
			}
			break
		}

		if msgType != websocket.TextMessage {
			continue
		}

		// Parse top-level command
		if err = json.Unmarshal(message, &cmd); err != nil {
			continue
		}

		switch cmd.Type {
		case "subscribe":
			fmt.Println("got subscribe request!")
			var sub SubscribeCmd
			if json.Unmarshal(cmd.Data, &sub); sub.Name == "" {
				continue
			}
			addr64, err := strconv.ParseUint(sub.Addr, 10, 64)
			if err != nil {
				continue
			}

			MainSimulation.AddDataSubscriber(
				sub.Name,
				ledger.Addr(addr64),
				sub.Etype,
				user.TraderID,
				conn,
			)

		case "unsubscribe":
			var unsub UnsubscribeCmd
			if json.Unmarshal(cmd.Data, &unsub); unsub.Addr == 0 {
				continue
			}
			MainSimulation.RemoveDataSubscriber(
				unsub.Addr,
				unsub.Etype,
				user.TraderID,
			)

		case "search":
			var search SearchCmd
			if json.Unmarshal(cmd.Data, &search); search.Name == "" {
				continue
			}
			results := MainSimulation.SearchDataSources(search.Name)
			_ = conn.WriteJSON(map[string]interface{}{
				"type": "search_results",
				"data": results,
			})

		case "swap":
			var swap SwapCmd
			if json.Unmarshal(cmd.Data, &swap); swap.AmountIn == 0 {
				continue
			}
			MainSimulation.PlaceUserSwap(
				user.TraderID,
				swap.FromSymbol,
				swap.ToSymbol,
				swap.AmountIn,
			)
			conn.WriteJSON(map[string]string{"status": "swap_placed"})
		}
	}
}