package ledger

import (
	"time"
	"sync"
	"fmt"
	"github.com/gorilla/websocket"
)

const data_source_timeout = 15 * time.Second

type data_source struct {
	Name		string
	Addr		Addr			// address of entity being watched
	Etype		EntityType
	Data		chan map[string]interface{}
	Cancel		chan struct{}	// to stop the runner
	mu			sync.RWMutex
	Subscribers map[uint64]*websocket.Conn
}

func new_data_source(name string, addr Addr, etype EntityType) *data_source {
	return &data_source{
		Name: name,
		Addr: addr,
		Etype: etype,
		Data: make(chan map[string]interface{}, 100),
		Cancel: make(chan struct{}),
		Subscribers: make(map[uint64]*websocket.Conn),
	}
}

func (d *data_source) AddSubscriber(id uint64, conn *websocket.Conn) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.Subscribers[id] = conn
	fmt.Println("\n\n\n\n\n\nadded subscriber")

}

func (d *data_source) RemoveSubscriber(id uint64) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.Subscribers[id]; ok {
		delete(d.Subscribers, id)
	}
}

func (d *data_source) Emit(tick uint64, data interface{}, etype EntityType) {
	payload := map[string]interface{}{
		"src":  d.Name,
		"type": etype.String(),
		"data": data,
		"tick": tick,
	}

	select {
	case d.Data <- payload:
	case <-time.After(1 * time.Second):
		fmt.Printf("[ %s ] emit channel full, dropping message\n", d.Name)
	}
}

func (d *data_source) broadcast(msg map[string]interface{}) {
	var to_remove []uint64

	to_remove = make([]uint64, 0)
	d.mu.RLock()
	defer d.mu.RUnlock()

	for id, conn := range d.Subscribers {
		if err := conn.WriteJSON(msg); err != nil {
			fmt.Printf("[ %s ] failed writing to subscriber %d → removing: %v\n", d.Name, id, err)
			to_remove = append(to_remove, id)
		}
	}

	for _, id := range to_remove {
		delete(d.Subscribers, id)
	}
}

func (d *data_source) ping() {
	payload := map[string]interface{}{
		"src":  d.Name,
		"type": "ping",
		"ts":   time.Now().Unix(),
	}
	d.broadcast(payload)
}

func (d *data_source) Run() {
	for {
		select {
		case <-d.Cancel:
			fmt.Printf("[ %s ] data_source stopped\n", d.Name)
			return
		case msg := <-d.Data:
			d.broadcast(msg)
		case <-time.After(5 * time.Second):
			d.ping()
		}
	}
}

func (d *data_source) Stop() {
	close(d.Cancel)
}