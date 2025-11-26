package ledger

import (
	"time"
	"sync"
	"fmt"
	"encoding/json"

	"github.com/gorilla/websocket"
)

const data_source_timeout = 15 * time.Second

type data_source struct {
	Name		string
	Addr		Addr			// address of entity being watched
	Etype		EntityType
	Data		chan []byte
	Cancel		chan struct{}	// to stop the runner
	mu			sync.RWMutex
	Subscribers map[uint64]*websocket.Conn
}

func new_data_source(name string, addr Addr, etype EntityType) *data_source {
	return &data_source{
		Name: name,
		Addr: addr,
		Etype: etype,
		Data: make(chan []byte, 100),
		Cancel: make(chan struct{}),
		Subscribers: make(map[uint64]*websocket.Conn),
	}
}

func (d *data_source) AddSubscriber(id uint64, conn *websocket.Conn) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.Subscribers[id] = conn
}

func (d *data_source) RemoveSubscriber(id uint64) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.Subscribers[id]; ok {
		delete(d.Subscribers, id)
	}
}

func (d *data_source) Emit(data interface{}, etype EntityType) {
	payload := map[string]interface{}{
		"src":  d.Name,
		"type": etype.String(),
		"data": data,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("[ %s ] failed to marshal emit data: %v\n", d.Name, err)
		return
	}

	select {
	case d.Data <- b:
	case <-time.After(1 * time.Second):
		fmt.Printf("[ %s ] emit channel full, dropping message\n", d.Name)
	}
}

func (d *data_source) broadcast(msg []byte) {
	var to_remove []uint64

	to_remove = make([]uint64, 0)
	d.mu.RLock()
	defer d.mu.RUnlock()

	for id, conn := range d.Subscribers {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			fmt.Printf("[ %s ] failed writing to subscriber %d → removing: %v\n", d.Name, id, err)
			to_remove = append(to_remove, id)
		}
	}

	for id := range to_remove {
		delete(d.Subscribers, id)
	}
}

func (d *data_source) ping() {
	payload := map[string]interface{}{
		"src":  d.Name,
		"type": "ping",
		"ts":   time.Now().Unix(),
	}
	b, _ := json.Marshal(payload)
	d.broadcast(b)
}

func (d *data_source) Run() {
	ticker := time.NewTicker(data_source_timeout)
	defer ticker.Stop()

	for {
		select {
		case <-d.Cancel:
			fmt.Printf("[ %s ] data_source stopped\n", d.Name)
			return
		case msg := <-d.Data:
			d.broadcast(msg)
		case <-ticker.C:
			d.ping()
		}
	}
}

func (d *data_source) Stop() {
	close(d.Cancel)
}