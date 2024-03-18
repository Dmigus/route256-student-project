package ordersgetter

import (
	"context"
	"route256.ozon.ru/project/loms/internal/models"
)

type ordersStorage interface {
	Load(context.Context, int64) (*models.Order, error)
}

type OrdersGetter struct {
	orders ordersStorage
}

func (og *OrdersGetter) Get(ctx context.Context, orderId int64) (*models.Order, error) {
	return og.orders.Load(ctx, orderId)
}
