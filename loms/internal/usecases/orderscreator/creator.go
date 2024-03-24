package orderscreator

import (
	"context"
	"errors"
	"fmt"
	"route256.ozon.ru/project/loms/internal/models"
)

type orderIdGenerator interface {
	NewId() int64
}

type stocksStorage interface {
	Reserve(context.Context, []models.OrderItem) error
}

type ordersStorage interface {
	Save(context.Context, *models.Order) error
}

type OrdersCreator struct {
	orderIdGenerator orderIdGenerator
	orders           ordersStorage
	stocks           stocksStorage
}

func NewOrdersCreator(orderIdGenerator orderIdGenerator, orders ordersStorage, stocks stocksStorage) *OrdersCreator {
	return &OrdersCreator{orderIdGenerator: orderIdGenerator, orders: orders, stocks: stocks}
}

func (oc *OrdersCreator) Create(ctx context.Context, userId int64, items []models.OrderItem) (int64, error) {
	newOrder := oc.createOrderInstance(userId)
	errReserving := oc.stocks.Reserve(ctx, items)
	if errReserving != nil {
		errReserving = fmt.Errorf("could not reserve items for user %d: %w", userId, errReserving)
		newOrder.Status = models.Failed
	} else {
		newOrder.Status = models.AwaitingPayment
		newOrder.IsItemsReserved = true
	}
	newOrder.Items = items
	errSaving := oc.orders.Save(ctx, newOrder)
	if errSaving != nil {
		errSaving = fmt.Errorf("could not save created order for user %d: %w", userId, errSaving)
	}
	errs := errors.Join(errSaving, errReserving)
	if errs != nil {
		return 0, errs
	}
	return newOrder.Id(), nil
}

func (oc *OrdersCreator) createOrderInstance(userId int64) *models.Order {
	newOrderId := oc.orderIdGenerator.NewId()
	return models.NewOrder(userId, newOrderId)
}
