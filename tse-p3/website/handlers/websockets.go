package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"tse-p3/globals"
	"tse-p3/ledger"
	"tse-p3/simulation"
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
	Etype ledger.EntityType `json:"etype"`
	Addr  ledger.Addr       `json:"addr"`
}

type UnsubscribeCmd struct {
	Etype ledger.EntityType `json:"etype"`
	Addr  ledger.Addr       `json:"addr"`
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
		conn     *websocket.Conn
		err      error
		ctx      context.Context
		user     users.User
		userOK   bool
		userVal  interface{}
		nameVal  interface{}
		subsVal  interface{}
		name     string
		subs     []users.DataSubscription
		msgType  int
		message  []byte
		cmd      CommandMessage
	)

	// 1. Upgrade connection
	conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("WebSocket upgrade failed: %v\n", err)
		return
	}
	defer conn.Close()

	// 2. Extract user from context (set by auth middleware or login/register)
	ctx = r.Context()
	userVal = ctx.Value("user")
	if userVal != nil {
		user, userOK = userVal.(users.User)
	}
	if !userOK || user.ID == 0 {
		conn.WriteJSON(map[string]string{"error": "unauthenticated"})
		return
	}

	// Optional: also read from context keys if you only stored name + subs
	if nameVal = ctx.Value("user.name"); nameVal != nil {
		name, _ = nameVal.(string)
	}
	if subsVal = ctx.Value("user.subscriptions"); subsVal != nil {
		subs, _ = subsVal.([]users.DataSubscription)
	}
	if name == "" {
		name = user.Name
	}
	if subs == nil {
		subs = user.DataSubscriptions
	}

	// 3. Re-subscribe to all saved data subscriptions on connect
	for _, sub := range subs {
		MainSimulation.AddDataSubscriber(
			sub.Name,
			sub.Addr,
			sub.Etype,
			user.TraderID,
			conn,
		)
	}

	// Send initial welcome / candle dump if you want
	_ = conn.WriteJSON(map[string]interface{}{
		"src":  "server",
		"type": "welcome",
		"data": map[string]string{"message": "connected as " + name},
	})

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
			var sub SubscribeCmd
			if json.Unmarshal(cmd.Data, &sub); sub.Name == "" {
				continue
			}
			MainSimulation.AddDataSubscriber(
				sub.Name,
				sub.Addr,
				sub.Etype,
				user.TraderID,
				conn,
			)

			// Persist subscription
			user.DataSubscriptions = append(user.DataSubscriptions, users.DataSubscription{
				Name:  sub.Name,
				Etype: sub.Etype,
				Addr:  sub.Addr,
			})
			_ = user.UpdateSubscriptions(r.Context())

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

			// Remove from user's saved subscriptions
			for i, s := range user.DataSubscriptions {
				if s.Addr == unsub.Addr && s.Etype == unsub.Etype {
					user.DataSubscriptions = append(user.DataSubscriptions[:i], user.DataSubscriptions[i+1:]...)
					break
				}
			}
			_ = user.UpdateSubscriptions(r.Context())

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
			err = MainSimulation.PlaceUserSwap(
				uint64(user.TraderID),
				swap.FromSymbol,
				swap.ToSymbol,
				swap.AmountIn,
			)
			if err != nil {
				conn.WriteJSON(map[string]string{"error": err.Error()})
			} else {
				conn.WriteJSON(map[string]string{"status": "swap_placed"})
			}
		}
	}
	// --- Remove Connection from all data sources --- //
	for _, dsub := range subs {
		MainSimulation.RemoveDataSubscriber(
			dsub.Addr,
			dsub.Etype,
			user.TraderID,
		)
	}
}