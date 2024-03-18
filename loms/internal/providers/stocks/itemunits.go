package stocks

import (
	"fmt"
	"sync"
)

var errNotEnoughItems = fmt.Errorf("required number of units is not available")

type itemUnits struct {
	mu              sync.RWMutex
	total, reserved uint64
}

func (g *itemUnits) getNumOfAvailable() uint64 {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.total - g.reserved
}

func (g *itemUnits) reserve(count uint16) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	available := g.total - g.reserved
	required := uint64(count)
	if available < required {
		return errNotEnoughItems
	}
	g.reserved += required
	return nil
}

func (g *itemUnits) cancelReserve(count uint16) {
	g.mu.Lock()
	defer g.mu.Unlock()
	freeNum := uint64(count)
	g.reserved -= freeNum
}

func (g *itemUnits) removeReserved(count uint16) {
	g.mu.Lock()
	defer g.mu.Unlock()
	freeNum := uint64(count)
	g.reserved -= freeNum
	g.total -= freeNum
}
