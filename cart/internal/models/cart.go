package models

import (
	"context"
	"sync"
)

type Cart struct {
	mu    sync.Mutex
	items map[int64]uint16
}

func NewCart() *Cart {
	return &Cart{
		items: make(map[int64]uint16),
	}
}

func (u *Cart) Add(_ context.Context, skuId int64, count uint16) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.items[skuId] += count
}

func (u *Cart) Delete(_ context.Context, skuId int64) {
	u.mu.Lock()
	defer u.mu.Unlock()
	delete(u.items, skuId)
}

func (u *Cart) Clear(_ context.Context) {
	u.mu.Lock()
	defer u.mu.Unlock()
	clear(u.items)
}

func (u *Cart) ListItems(_ context.Context) []CartItem {
	u.mu.Lock()
	defer u.mu.Unlock()
	items := make([]CartItem, 0, len(u.items))
	for skuId, count := range u.items {
		items = append(items, CartItem{SkuId: skuId, Count: count})
	}
	return items
}
