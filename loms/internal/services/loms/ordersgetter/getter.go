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
		WithinTransaction(context.Context, func(ctx context.Context, orders OrderRepo) bool) error
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
	var businessErr error
	trErr := og.tx.WithinTransaction(ctx, func(ctx context.Context, orders OrderRepo) bool {
		order, businessErr = orders.Load(ctx, orderID)
		if businessErr != nil {
			businessErr = fmt.Errorf("could not load order %d: %w", orderID, businessErr)
			return false
		}
		return true
	})
	if businessErr != nil {
		return nil, businessErr
	}
	if trErr != nil {
		return nil, trErr
	}
	return order, nil
}
