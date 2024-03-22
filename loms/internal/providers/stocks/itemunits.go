package stocks

import (
	"errors"
	"sync"
)

var ErrNotEnoughItems = errors.New("required number of units is not available")

type ItemUnits struct {
	mu              sync.RWMutex
	total, reserved uint64
}

func NewItemUnits(total, reserved uint64) *ItemUnits {
	return &ItemUnits{
		total:    total,
		reserved: reserved,
	}
}

func (g *ItemUnits) getNumOfAvailable() uint64 {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.total - g.reserved
}

func (g *ItemUnits) reserve(count uint16) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	available := g.total - g.reserved
	required := uint64(count)
	if available < required {
		return ErrNotEnoughItems
	}
	g.reserved += required
	return nil
}

func (g *ItemUnits) cancelReserve(count uint16) {
	g.mu.Lock()
	defer g.mu.Unlock()
	freeNum := uint64(count)
	g.reserved -= freeNum
}

func (g *ItemUnits) removeReserved(count uint16) {
	g.mu.Lock()
	defer g.mu.Unlock()
	freeNum := uint64(count)
	g.reserved -= freeNum
	g.total -= freeNum
}
