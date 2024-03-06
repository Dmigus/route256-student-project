package inmemorycart

import (
	"context"
	"route256.ozon.ru/project/cart/internal/usecases"
	"route256.ozon.ru/project/cart/internal/usecases/lister"
	"sync"
)

type InMemoryCart struct {
	items map[usecases.SkuId]usecases.ItemCount
	mu    sync.Mutex
}

func New() *InMemoryCart {
	return &InMemoryCart{
		items: make(map[usecases.SkuId]usecases.ItemCount),
	}
}

func (u *InMemoryCart) Add(_ context.Context, skuId usecases.SkuId, count usecases.ItemCount) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.items[skuId] += count
	return nil
}

func (u *InMemoryCart) Delete(_ context.Context, skuId usecases.SkuId) error {
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

func (u *InMemoryCart) ListItems(_ context.Context) ([]lister.CartItem, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	items := make([]lister.CartItem, 0, len(u.items))
	for skuId, count := range u.items {
		items = append(items, lister.CartItem{SkuId: skuId, Count: count})
	}
	return items, nil
}
