package ordersgetter

import (
	"context"
	"fmt"
	"route256.ozon.ru/project/loms/internal/models"
)

type ordersStorage interface {
	Load(context.Context, int64) (*models.Order, error)
}

type OrdersGetter struct {
	orders ordersStorage
}

func NewOrdersGetter(orders ordersStorage) *OrdersGetter {
	return &OrdersGetter{orders: orders}
}

func (og *OrdersGetter) Get(ctx context.Context, orderId int64) (*models.Order, error) {
	loaded, err := og.orders.Load(ctx, orderId)
	if err != nil {
		return nil, fmt.Errorf("could not load order %d: %w", orderId, err)
	}
	return loaded, nil
}
