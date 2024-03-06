package models

import (
	"context"
	"sync"
)

type InMemoryCart struct {
	items map[SkuId]ItemCount
	mu    sync.Mutex
}

func NewInMemoryCart() *InMemoryCart {
	return &InMemoryCart{
		items: make(map[SkuId]ItemCount),
	}
}

func (u *InMemoryCart) Add(_ context.Context, skuId SkuId, count ItemCount) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.items[skuId] += count
	return nil
}

func (u *InMemoryCart) Delete(_ context.Context, skuId SkuId) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	delete(u.items, skuId)
	return nil
}

func (u *InMemoryCart) Clear(_ context.Context) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	clear(u.items)
	return nil
}

func (u *InMemoryCart) ListItems(_ context.Context) ([]CartItem, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	items := make([]CartItem, 0, len(u.items))
	for skuId, count := range u.items {
		items = append(items, CartItem{SkuId: skuId, Count: count})
	}
	return items, nil
}
