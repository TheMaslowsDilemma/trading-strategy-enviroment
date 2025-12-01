package handlers

import (
	"fmt"
	"context"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	var (
		client		*ws_client
		conn		*websocket.Conn
		ctx			context.Context
		msg_type	int
		msg_in		[]byte
		err			error
	)

	ctx = r.Context()
	conn, err = upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Printf("websocket upgrade failed: %v\n", err)
		return
	}
	defer conn.Close()


	client, err = init_ws_client(ctx, conn)
	if err != nil {
		fmt.Printf("ws client init failed: %v\n", err)
		return
	}

	go client.WriterLoop()

	message_loop:
	for {
		msg_type, msg_in, err = conn.ReadMessage()
		if err != nil {
			if closeErr, ok := err.(*websocket.CloseError); ok {
				switch closeErr.Code {
					case websocket.CloseGoingAway:
						fmt.Printf("websocket closed normally (Going Away): %v\n", err)
						break message_loop

					case websocket.CloseNormalClosure:
						fmt.Printf("websocket closed normally: %v\n", err)
						break message_loop
				}
				break message_loop
			}

			// If we get here, it's either an unexpected close code or a non-CloseError
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure,websocket.CloseAbnormalClosure) {
				fmt.Printf("websocket read message error (unexpected): '%v'. stopping message loop.\n", err)
			} else {
				fmt.Printf("some other error occurred: %v\n", err)
			}

			break message_loop
		}
		
		err = client.handle_incoming(msg_type, msg_in)
		if err != nil {
			fmt.Printf("client failed to handle incoming: %v\n", err)
		}
	}
	client.Close()
}