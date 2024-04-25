// Package orderscreator содержит логику работы юзкейса создания заказа
package orderscreator

import (
	"context"
	"fmt"
	anotherErrors "github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"route256.ozon.ru/project/loms/internal/models"
)

var tracer = otel.Tracer("order creation")

// ErrInsufficientStocks это ошибка, обозначающая нехватку стоков для осуществления заказа
var ErrInsufficientStocks = anotherErrors.Wrap(models.ErrFailedPrecondition, "insufficient stocks")

type (
	// StockRepo это контракт для использования репозитория стоков OrdersCreator'ом. Используется другими слоями для настройки атомарности
	StockRepo interface {
		Reserve(context.Context, []models.OrderItem) error
	}
	// OrderRepo это контракт для использования репозитория заказов OrdersCreator'ом. Используется другими слоями для настройки атомарности
	OrderRepo interface {
		Create(context.Context, int64, []models.OrderItem) (*models.Order, error)
		Save(context.Context, *models.Order) error
	}
	// EventSender это контракт для использования отправителя уведомлений о изменении статуса заказа
	EventSender interface {
		OrderStatusChanged(context.Context, *models.Order) error
	}
	txManager interface {
		WithinTransaction(context.Context, func(ctx context.Context, orders OrderRepo, stocks StockRepo, evSender EventSender) bool) error
	}
	// OrdersCreator - сущность, которая умеет создавать заказы в системе
	OrdersCreator struct {
		tx txManager
	}
)

// NewOrdersCreator создаёт OrdersCreator. tx - должен быть объектом, позволяющим исполнять функцию атомарно
func NewOrdersCreator(tx txManager) *OrdersCreator {
	return &OrdersCreator{tx: tx}
}

// Create создаёт заказ для пользователя userID и товарами items и возращает его id. Атомарность всей операции обеспечивается объектом tx, переданым при создании
func (oc *OrdersCreator) Create(ctx context.Context, userID int64, items []models.OrderItem) (orderID int64, err error) {
	var businessErr error
	trErr := oc.tx.WithinTransaction(ctx, func(ctx context.Context, orders OrderRepo, stocks StockRepo, evSender EventSender) bool {
		orderID, businessErr = createOrder(ctx, userID, items, orders, stocks, evSender)
		if businessErr != nil && !anotherErrors.Is(businessErr, ErrInsufficientStocks) {
			return false
		}
		return true
	})
	if businessErr != nil {
		return 0, businessErr
	}
	if trErr != nil {
		return 0, trErr
	}
	return orderID, nil
}

func createOrder(ctx context.Context, userID int64, items []models.OrderItem, orders OrderRepo, stocks StockRepo, evSender EventSender) (_ int64, err error) {
	ctx, span := tracer.Start(ctx, "creation")
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()
	order, err := registerNewOrder(ctx, userID, items, orders, evSender)
	if err != nil {
		return 0, err
	}
	span.AddEvent("order registered in system")
	errReserving := assignStocksForOrder(ctx, order, stocks)
	if err = saveFinalOrderState(ctx, order, orders, evSender); err != nil {
		return 0, err
	}
	if errReserving != nil {
		return 0, errReserving
	}
	span.AddEvent("item stock was reserved")
	return order.Id(), nil
}

func registerNewOrder(ctx context.Context, userID int64, items []models.OrderItem, orders OrderRepo, evSender EventSender) (*models.Order, error) {
	order, err := orders.Create(ctx, userID, items)
	if err != nil {
		return nil, fmt.Errorf("could not create new order for user %d: %w", userID, err)
	}
	if err = evSender.OrderStatusChanged(ctx, order); err != nil {
		return nil, fmt.Errorf("could not order changing status with id = %d: %w", order.Id(), err)
	}
	return order, nil
}

func assignStocksForOrder(ctx context.Context, order *models.Order, stocks StockRepo) error {
	err := stocks.Reserve(ctx, order.Items)
	if err != nil {
		err = fmt.Errorf("could not reserve items for order %d: %w", order.Id(), err)
		order.Status = models.Failed
		return err
	}
	order.Status = models.AwaitingPayment
	order.IsItemsReserved = true
	return nil
}

func saveFinalOrderState(ctx context.Context, order *models.Order, orders OrderRepo, evSender EventSender) error {
	if err := orders.Save(ctx, order); err != nil {
		return fmt.Errorf("could not save created order with id = %d: %w", order.Id(), err)
	}
	if err := evSender.OrderStatusChanged(ctx, order); err != nil {
		return fmt.Errorf("could not send order changing status with id = %d: %w", order.Id(), err)
	}
	return nil
}
