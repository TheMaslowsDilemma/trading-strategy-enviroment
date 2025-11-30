package handlers

import (
	"fmt"
	"context"
	"strconv"
	"encoding/json"

	"tse-p3/users"
	"tse-p3/ledger"
	"github.com/gorilla/websocket"
)

const WS_CLIENT_OUT_BUFF_SIZE = 100

type ws_client struct {
	conn	*websocket.Conn
	ctx		context.Context
	Healthy	bool
	user	users.User
	Cancel	chan bool
	MsgOut	chan map[string]interface{}
}

func (client *ws_client) Close() {
	var err error

	if client == nil {
		return
	}
	client.Healthy = false
	client.Cancel <- true

	err = users.SetUserActivity(client.ctx, client.user.ID, false)
	if err != nil {
		fmt.Printf("ws client failed to update user '%v' activity: %v", client.user.ID, err)
	}


	// tell the simulation to 
	// 1. save updates to wallet amounts for the client into pg
	// 2. remove wallet, trader, subscribers from simulation
}

func init_ws_client(ctx context.Context, wsc *websocket.Conn) (*ws_client, error) {
	var (
		usr	users.User
		ok	bool
		err	error
	)

	if v := ctx.Value("user"); v != nil {
		usr, ok = v.(users.User)
		if !ok {
			return nil, fmt.Errorf("context's 'user' could not be cast to User")
		}
	} else {
		return nil, fmt.Errorf("context has no 'user'")
	}

	if !usr.Active {
		err = users.SetUserActivity(ctx, usr.ID, true)
		if err != nil {
			return nil, fmt.Errorf("failed to update user activity: %v", err)
		}

		// TODO:
		// go and turn on its trader / wallets etc.
		// this means those wallets etc. would need to
		// be persisted somewhere so that at this point
		// they could start with the same balances
	}


	return &ws_client {
		conn: wsc,
		ctx: ctx,
		Healthy: true,
		user: usr,
		Cancel: make(chan bool),
		MsgOut:  make(chan map[string]interface{}, WS_CLIENT_OUT_BUFF_SIZE),
	}, nil
}

func (client *ws_client) emit(msg map[string]interface{}) error {
	if client == nil {
		return fmt.Errorf("client is nil")
	}
	
	if !client.Healthy {
		return fmt.Errorf("client '%v' is not healthy", client.user.Name)
	}

	select {
	case client.MsgOut <- msg:
	default:
		return fmt.Errorf("client '%v' has a full msg buffer", client.user.Name)
	}
	return nil
}


func (client *ws_client) write(msg map[string]interface{}) error {
	var err error
	err = client.conn.WriteJSON(msg)
	if err != nil {
		return fmt.Errorf("failed to write json to conn: %v", err)
	}
	return nil
}

func (client *ws_client) WriterLoop() {
	var outmsg map[string]interface{}
	var ok	bool
	for {
		select {
			case _ = <-client.Cancel:
				fmt.Printf("stopping client '%v' writer loop", client.user.Name)
				return
			case outmsg, ok = <-client.MsgOut:
				if !ok {
					continue
				}
				client.write(outmsg)
		}
	}
}

func (client *ws_client) handle_incoming(msg_type int, msg_in []byte) error {
	var (
		err error
		cmd	CommandMessage
	)

	if msg_type != websocket.TextMessage {
		return fmt.Errorf("non text message")
	}

	if err = json.Unmarshal(msg_in, &cmd); err != nil {
		return fmt.Errorf("failed to marshal message as 'command': %v", err)
	}

	switch cmd.Type {
	case "subscribe":
		err = client.cmd_subscribe(cmd)
		if err != nil {
			return fmt.Errorf("client '%v' failed to handle subscribe: %v", client.user.Name, err)
		}
	case "unsubscribe":
		err = client.cmd_unsubscribe(cmd)
		if err != nil {
			return fmt.Errorf("client '%v' failed to handle unsubscribe: %v", client.user.Name, err)
		}
	case "search":
		fmt.Println("searching")
		err = client.cmd_search(cmd)
		if err != nil {
			return fmt.Errorf("client '%v' failed to handle search: %v", client.user.Name, err)
		}
	case "swap":
		err = client.cmd_swap(cmd)
		if err != nil {
			return fmt.Errorf("client '%v' failed to handle swap: %v", client.user.Name, err)
		}
	default:
		return fmt.Errorf("unknown cmd type '%v'\n", cmd.Type)
	}
	return nil
}

func (client *ws_client) cmd_search(cmd CommandMessage) error {
	var (
		srch	SearchCmd
		rslt	[]ledger.SearchResult
		err		error
	)

	if json.Unmarshal(cmd.Data, &srch); srch.Name == "" {
		return fmt.Errorf("no search 'name' given")
	}

	rslt = MainSimulation.SearchDataSources(srch.Name)
	err = client.emit(map[string] interface{} {
		"type": "search_results",
		"data": rslt,
	})
	fmt.Printf("results: %v\n", rslt)
	if err != nil {
		return fmt.Errorf("failed enqueue for emission: %v", err)
	}
	return nil
}


func (client *ws_client) cmd_subscribe(cmd CommandMessage) error {
	var (
		sub		SubscribeCmd
		uaddr	uint64
		err		error
	)

	if json.Unmarshal(cmd.Data, &sub); sub.Name == "" {
		return fmt.Errorf("failed to unmarshal subscribe cmd")
	}

	uaddr, err = strconv.ParseUint(sub.Addr, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse address: %v", err)
	}

	MainSimulation.AddDataSubscriber(
		sub.Name,
		ledger.Addr(uaddr),
		sub.Etype,
		client.user.TraderID,
		client.emit,
	)
	return nil
}


func (client *ws_client) cmd_unsubscribe(cmd CommandMessage) error {
	var (
		unsub	UnsubscribeCmd
		uaddr	uint64
		err		error
	)

	if json.Unmarshal(cmd.Data, &unsub); unsub.Addr == "" {
		return fmt.Errorf("failed to unmarshal unsubscribe cmd")
	}

	uaddr, err = strconv.ParseUint(unsub.Addr, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse address: %v", err)
	}

	MainSimulation.RemoveDataSubscriber(
		ledger.Addr(uaddr),
		unsub.Etype,
		client.user.TraderID,
	)
	return nil
}

func (client *ws_client) cmd_swap(cmd CommandMessage) error {
	var (
		swap	SwapCmd
	)
	
	if json.Unmarshal(cmd.Data, &swap); swap.FromSymbol == "" {
		return fmt.Errorf("failed to unmarshal swap cmd")
	}
	
	MainSimulation.PlaceUserSwap(
		client.user.TraderID,
		swap.FromSymbol,
		swap.ToSymbol,
		swap.AmountIn,
	)
	return nil
}


