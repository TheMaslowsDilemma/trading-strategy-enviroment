package ledger

import (
	"fmt"
	"sync"
)

/** Question: Why do data_source need their own go routines... **/

// ------------------------ EmitterManager ---------------------------- //

type Emit func (msg map[string]interface{}) error

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
	case Wallet_t:
		m = em.Wallets
	case Exchange_t:
		m = em.Exchanges
	default:
		return nil
	}

	if dsrc, exists := m[addr]; exists {
		return dsrc // already exists
	}

	dsrc := new_data_source(name, addr, etype)
	m[addr] = dsrc
	em.DataCatalog.AddSource(dsrc)

	go dsrc.Run() // start the data source routine

	return dsrc
}

func (em *EmitterManager) RemoveSource(addr Addr, etype EntityType) {
	var (
		ds	*data_source
		ok	bool
		m 	map[Addr]*data_source
	)

	em.mu.Lock()
	defer em.mu.Unlock()

	switch etype {
	case Wallet_t:
		m = em.Wallets
	case Exchange_t:
		m = em.Exchanges
	default:
		return
	}
	
	if ds, ok = m[addr]; ok {
		ds.Stop()
		em.DataCatalog.RemoveSource(ds)
		delete(m, addr)
	}
}

// AddSubscriber attaches a WebSocket connection to a source (creates source if needed)
func (em *EmitterManager) AddSubscriber(name string, addr Addr, etype EntityType, userID uint64, emitter Emit) {
	dsrc := em.AddSource(name, addr, etype)
	dsrc.AddEmitter(userID, emitter)
	fmt.Printf("Subscriber '%v' connected to '%v' %v\n", userID, name, etype.String())
}

func (em *EmitterManager) RemoveSubscriber(addr Addr, etype EntityType, userID uint64) {
	em.mu.RLock()
	defer em.mu.RUnlock()

	var m map[Addr]*data_source
	switch etype {
	case Wallet_t:
		m = em.Wallets
	case Exchange_t:
		m = em.Exchanges
	}

	if src, ok := m[addr]; ok {
		src.RemoveEmitter(userID)
	}
}