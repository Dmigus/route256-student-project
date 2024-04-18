// Package orderscanceller содержит логику работы юзкейса отмены заказа
package orderscanceller

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"route256.ozon.ru/project/loms/internal/models"
)

var tracer = otel.Tracer("order cancelling")

var errWrongOrderStatus = errors.Wrap(models.ErrFailedPrecondition, "order status is wrong")

type (
	// OrderRepo это контракт для использования репозитория заказов OrderCanceller'ом. Используется другими слоями для настройки атомарности
	OrderRepo interface {
		Save(context.Context, *models.Order) error
		Load(context.Context, int64) (*models.Order, error)
	}
	// StockRepo это контракт для использования репозитория стоков OrderCanceller'ом. Используется другими слоями для настройки атомарности
	StockRepo interface {
		CancelReserved(context.Context, []models.OrderItem) error
		AddItems(context.Context, []models.OrderItem) error
	}
	// EventSender это контракт для использования отправителя уведомлений о изменении статуса заказа
	EventSender interface {
		OrderStatusChanged(context.Context, *models.Order) error
	}
	txManager interface {
		WithinTransaction(context.Context, func(ctx context.Context, orders OrderRepo, stocks StockRepo, evSender EventSender) bool) error
	}
	// OrderCanceller - сущность, которая умеет отменять заказы
	OrderCanceller struct {
		tx txManager
	}
)

// NewOrderCanceller создаёт OrderCanceller. tx - должен быть объектом, позволяющим исполнять функцию атомарно
func NewOrderCanceller(tx txManager) *OrderCanceller {
	return &OrderCanceller{tx: tx}
}

// Cancel отменяет заказ с id = orderId. Атомарность всей операции обеспечивается объектом tx, переданым при создании
func (oc *OrderCanceller) Cancel(ctx context.Context, orderID int64) error {
	var businessErr error
	trErr := oc.tx.WithinTransaction(ctx, func(ctx context.Context, orders OrderRepo, stocks StockRepo, evSender EventSender) bool {
		businessErr = cancelOrder(ctx, orderID, orders, stocks, evSender)
		return businessErr == nil
	})
	if businessErr != nil {
		return businessErr
	}
	return trErr
}

func cancelOrder(ctx context.Context, orderID int64, orders OrderRepo, stocks StockRepo, evSender EventSender) (err error) {
	ctx, span := tracer.Start(ctx, "cancelling")
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()
	order, err := orders.Load(ctx, orderID)
	if err != nil {
		err = fmt.Errorf("could not load order %d: %w", orderID, err)
		return err
	}
	if err = cancelAnyOrder(ctx, stocks, order); err != nil {
		return err
	}
	span.AddEvent("order cancelled")
	if err = saveFinalOrderState(ctx, order, orders, evSender); err != nil {
		return err
	}
	return nil
}

func cancelAnyOrder(ctx context.Context, stocks StockRepo, order *models.Order) error {
	if order.Status == models.Cancelled {
		return errWrongOrderStatus
	}
	if order.IsItemsReserved {
		if err := cancelReserved(ctx, stocks, order); err != nil {
			return err
		}
	} else if order.Status == models.Payed {
		if err := cancelPayed(ctx, stocks, order); err != nil {
			return err
		}
	} else {
		order.Status = models.Cancelled
	}
	return nil
}

func cancelReserved(ctx context.Context, stocks StockRepo, order *models.Order) error {
	if order.IsItemsReserved {
		err := stocks.CancelReserved(ctx, order.Items)
		if err != nil {
			return fmt.Errorf("could not cancel reserved items for order %d: %w", order.Id(), err)
		}
		order.IsItemsReserved = false
		order.Status = models.Cancelled
	}
	return nil
}

func cancelPayed(ctx context.Context, stocks StockRepo, order *models.Order) error {
	if order.Status == models.Payed {
		err := stocks.AddItems(ctx, order.Items)
		if err != nil {
			return fmt.Errorf("could not return payed items for order %d: %w", order.Id(), err)
		}
		order.Status = models.Cancelled
	}
	return nil
}

func saveFinalOrderState(ctx context.Context, order *models.Order, orders OrderRepo, evSender EventSender) error {
	if err := orders.Save(ctx, order); err != nil {
		return fmt.Errorf("could not save order with id = %d: %w", order.Id(), err)
	}
	if err := evSender.OrderStatusChanged(ctx, order); err != nil {
		return fmt.Errorf("could not send order changing status with id = %d: %w", order.Id(), err)
	}
	return nil
}
