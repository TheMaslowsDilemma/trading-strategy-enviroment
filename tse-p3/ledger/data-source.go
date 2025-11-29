package ledger

import (
	"time"
	"sync"
	"fmt"
)

// TODO:
// when a data-source is deleted, we must notify any active subscribers that it is gone.
// or we just let them figure that out because it doesn't update anymore...
// or we notify but in the UI just mark as [deleted]

const data_source_timeout = 15 * time.Second

type data_source struct {
	Name		string
	Addr		Addr			// address of entity being watched
	Etype		EntityType
	Data		chan map[string]interface{}
	Cancel		chan struct{}	// to stop the runner
	mu			sync.RWMutex
	Emissions map[uint64]Emit
}

func new_data_source(name string, addr Addr, etype EntityType) *data_source {
	return &data_source{
		Name: name,
		Addr: addr,
		Etype: etype,
		Data: make(chan map[string]interface{}, 100),
		Cancel: make(chan struct{}),
		Emissions: make(map[uint64]Emit),
	}
}

func (d *data_source) AddEmitter(id uint64, emitter Emit) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.Emissions[id] = emitter
}

func (d *data_source) RemoveEmitter(id uint64) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.Emissions[id]; ok {
		delete(d.Emissions, id)
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
	default:
		fmt.Printf("data source '%s' channel full, dropping message\n", d.Name)
	}
}

func (d *data_source) broadcast(msg map[string]interface{}) {
	var to_remove []uint64

	to_remove = make([]uint64, 0)
	d.mu.RLock()
	defer d.mu.RUnlock()

	for id, emitter := range d.Emissions {
		if err := emitter(msg); err != nil {
			fmt.Printf("data source '%s' failed emitting '%d': %v\n", d.Name, id, err)
			to_remove = append(to_remove, id)
		}
	}

	for _, id := range to_remove {
		delete(d.Emissions, id)
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
	ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
	for {
		select {
		case <-d.Cancel:
			fmt.Printf("data source '%s' stopped\n", d.Name)
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