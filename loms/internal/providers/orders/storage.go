package orders

import (
	"context"
	"fmt"
	"route256.ozon.ru/project/loms/internal/models"
	"sync"
)

var errOrderNotFound = fmt.Errorf("order not found")

type InMemoryOrdersStorage struct {
	mu   sync.RWMutex
	data map[int64]*models.Order
}

func (i *InMemoryOrdersStorage) Save(_ context.Context, order *models.Order) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.data[order.Id()] = order
	return nil
}

func (i *InMemoryOrdersStorage) Load(_ context.Context, orderId int64) (*models.Order, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	order, exists := i.data[orderId]
	if !exists {
		return nil, errOrderNotFound
	}
	return order, nil
}
