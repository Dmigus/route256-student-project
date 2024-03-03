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

func (u *InMemoryCart) Add(_ context.Context, skuId modifier.SkuId, count uint16) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.items[skuId] += count
	return nil
}

func (u *InMemoryCart) Delete(_ context.Context, skuId modifier.SkuId) error {
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

// Обходит все пары (ИД товара, количество) и вызывает функцию f до того момента, пока f возвращает true. Предполагается, что в случае ошибки f вернёт false и обход сразу завершится.
func (u *InMemoryCart) Range(ctx context.Context, f func(ctx context.Context, skuId modifier.SkuId, count uint16) bool) {
	u.mu.Lock()
	defer u.mu.Unlock()
	continueRange := true
	for skuId, count := range u.items {
		continueRange = f(ctx, skuId, count)
		if !continueRange {
			return
		}
	}
}
