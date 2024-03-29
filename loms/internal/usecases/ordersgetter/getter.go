package ordersgetter

import (
	"context"
	"fmt"
	"route256.ozon.ru/project/loms/internal/models"
)

type (
	OrderRepo interface {
		Load(context.Context, int64) (*models.Order, error)
	}
	txManager interface {
		WithinTransaction(context.Context, func(ctx context.Context, orders OrderRepo, _ any) error) error
	}
	OrdersGetter struct {
		tx txManager
	}
)

func NewOrdersGetter(tx txManager) *OrdersGetter {
	return &OrdersGetter{tx: tx}
}

func (og *OrdersGetter) Get(ctx context.Context, orderId int64) (order *models.Order, err error) {
	err = og.tx.WithinTransaction(ctx, func(ctx context.Context, orders OrderRepo, _ any) error {
		order, err = orders.Load(ctx, orderId)
		if err != nil {
			return fmt.Errorf("could not load order %d: %w", orderId, err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return order, nil
}
