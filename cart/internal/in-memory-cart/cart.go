package in_memory_cart

import (
	"context"
	"route256.ozon.ru/project/cart/internal/service/modifier"
	"sync"
)

type InMemoryCart struct {
	items map[modifier.SkuId]uint16 // мб стоит заменить на sync.Map
	mu    sync.Mutex
}

func New() *InMemoryCart {
	return &InMemoryCart{
		items: make(map[modifier.SkuId]uint16),
	}
}

func (u *InMemoryCart) Add(ctx context.Context, skuId modifier.SkuId, count uint16) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.items[skuId] += count
	return nil
}

func (u *InMemoryCart) Delete(ctx context.Context, skuId modifier.SkuId) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	delete(u.items, skuId)
	return nil
}

func (u *InMemoryCart) Clear(ctx context.Context) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	clear(u.items)
	return nil
}

func (u *InMemoryCart) Range(ctx context.Context, f func(ctx context.Context, skuId modifier.SkuId, count uint16)) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	for skuId, count := range u.items {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			f(ctx, skuId, count)
		}
	}
	return nil
}
