// Package ordersgetter содержит логику работы юзкейса получения информации о заказе
package ordersgetter

import (
	"context"
	"fmt"
	"route256.ozon.ru/project/loms/internal/models"
)

type (
	// OrderRepo это контракт для использования репозитория заказов OrdersGetter'ом. Используется другими слоями для настройки доступа к исключительно зафиксированным данным
	OrderRepo interface {
		Load(context.Context, int64) (*models.Order, error)
	}
	txManager interface {
		WithinTransaction(context.Context, func(ctx context.Context, orders OrderRepo) error) error
	}
	// OrdersGetter - сущность, которая умеет возращать информацию о заказах в системе
	OrdersGetter struct {
		tx txManager
	}
)

// NewOrdersGetter создаёт OrdersGetter. tx - должен быть объектом, позволяющим читать только зафиксированные данные
func NewOrdersGetter(tx txManager) *OrdersGetter {
	return &OrdersGetter{tx: tx}
}

// Get возвращает информацию о заказе с id = orderID
func (og *OrdersGetter) Get(ctx context.Context, orderID int64) (order *models.Order, err error) {
	err = og.tx.WithinTransaction(ctx, func(ctx context.Context, orders OrderRepo) error {
		order, err = orders.Load(ctx, orderID)
		if err != nil {
			return fmt.Errorf("could not load order %d: %w", orderID, err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return order, nil
}
