// Package orderspayer содержит логику работы юзкейса оплаты заказа
package orderspayer

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"route256.ozon.ru/project/loms/internal/models"
)

var errWrongOrderStatus = errors.Wrap(models.ErrFailedPrecondition, "order status is wrong")

type (
	// OrderRepo это контракт для использования репозитория заказов OrdersPayer'ом. Используется другими слоями для настройки атомарности
	OrderRepo interface {
		Save(context.Context, *models.Order) error
		Load(context.Context, int64) (*models.Order, error)
	}
	// StockRepo это контракт для использования репозитория стоков OrdersPayer'ом. Используется другими слоями для настройки атомарности
	StockRepo interface {
		RemoveReserved(context.Context, []models.OrderItem) error
	}
	txManager interface {
		WithinTransaction(context.Context, func(ctx context.Context, orders OrderRepo, stocks StockRepo) error) error
	}
	// OrdersPayer - сущность, которая умеет осуществлять оплату заказа в системе
	OrdersPayer struct {
		tx txManager
	}
)

// NewOrdersPayer создаёт OrdersPayer. tx - должен быть объектом, позволяющим выполнять действия атомарно
func NewOrdersPayer(tx txManager) *OrdersPayer {
	return &OrdersPayer{tx: tx}
}

// Pay осуществлять оплату заказа с id = orderID. Атомарность всей операции обеспечивается объектом tx, переданым при создании
func (or *OrdersPayer) Pay(ctx context.Context, orderID int64) error {
	return or.tx.WithinTransaction(ctx, func(ctx context.Context, orders OrderRepo, stocks StockRepo) error {
		order, err := orders.Load(ctx, orderID)
		if err != nil {
			return fmt.Errorf("could not load order %d: %w", orderID, err)
		}
		if order.Status != models.AwaitingPayment {
			return errWrongOrderStatus
		}
		err = stocks.RemoveReserved(ctx, order.Items)
		if err != nil {
			return fmt.Errorf("could not remove reserved items for order %d: %w", orderID, err)
		}
		order.Status = models.Payed
		order.IsItemsReserved = false
		if err = orders.Save(ctx, order); err != nil {
			return fmt.Errorf("could not save order %d: %w", orderID, err)
		}
		return nil
	})
}
