package ledger

import (
	"fmt"
	"time"
	"sync"
	"encoding/json"
	"github.com/gorilla/websocket"
)

/** Question: Why do data_source need their own go routines... **/

// ------------------------ EmitterManager ---------------------------- //

type EmitterManager struct {
	mu			sync.RWMutex
	Wallets		map[Addr]*data_source
	Exchanges	map[Addr]*data_source
	DataCatalog *data_catalog
}

func NewEmitterManager() *EmitterManager {
	return &EmitterManager{
		Wallets:	make(map[Addr]*data_source),
		Exchanges:	make(map[Addr]*data_source),
		DataCatalog: new_data_catalog(),
	}
}


func (em *EmitterManager) SearchSources(name string) []SearchResult {
	return em.DataCatalog.SearchK(name)
}

// AddSource creates a new data source if it doesn't exist and starts its goroutine
func (em *EmitterManager) AddSource(name string, addr Addr, etype EntityType) *data_source {
	em.mu.Lock()
	defer em.mu.Unlock()

	var m map[Addr]*data_source
	switch etype {
	case EntityWallet:
		m = em.Wallets
	case EntityExchange:
		m = em.Exchanges
	default:
		return nil
	}

	if dsrc, exists := m[addr]; exists {
		return dsrc // already exists
	}

	dsrc := new_data_source(name, addr, etype)
	m[addr] = dsrc
	en.DataCatalog.AddSource(dsrc)

	go dsrc.Run() // start the data source routine

	return dsrc
}

// AddSubscriber attaches a WebSocket connection to a source (creates source if needed)
func (em *EmitterManager) AddSubscriber(name string, addr Addr, etype EntityType, userID uint64, conn *websocket.Conn) {
	dsrc := em.AddSource(name, addr, etype)
	dsrc.AddSubscriber(userID, conn)
	fmt.Printf("Subscriber %d connected to %s (%s)\n", userID, addr, etype.String())
}

func (em *EmitterManager) RemoveSubscriber(addr Addr, etype EntityType, userID uint64) {
	em.mu.RLock()
	defer em.mu.RUnlock()

	var m map[Addr]*data_source
	switch etype {
	case EntityWallet:
		m = em.Wallets
	case EntityExchange:
		m = em.Exchanges
	}

	if src, ok := m[addr]; ok {
		src.RemoveSubscriber(userID)
	}
}