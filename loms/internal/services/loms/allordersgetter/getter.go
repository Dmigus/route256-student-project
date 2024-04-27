// Package allordersgetter содержит реализацию юзкейса получения всех заказов
package allordersgetter

import (
	"context"
	"fmt"

	"route256.ozon.ru/project/loms/internal/models"
)

type (
	// OrderRepo это контракт для использования репозитория заказов OrdersGetter'ом. Используется другими слоями для настройки доступа к исключительно зафиксированным данным
	OrderRepo interface {
		LoadAll(context.Context) ([]*models.Order, error)
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

// Get возвращает информацию о всех заказах
func (og *OrdersGetter) Get(ctx context.Context) (orders []*models.Order, err error) {
	var businessErr error
	trErr := og.tx.WithinTransaction(ctx, func(ctx context.Context, orderRepo OrderRepo) bool {
		orders, businessErr = orderRepo.LoadAll(ctx)
		if businessErr != nil {
			businessErr = fmt.Errorf("could not load all orders: %w", businessErr)
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
	return orders, nil
}
