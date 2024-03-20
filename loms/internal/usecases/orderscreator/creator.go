package orderscreator

import (
	"context"
	"errors"
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
	if errReserving == nil {
		newOrder.Status = models.AwaitingPayment
	} else {
		newOrder.Status = models.Failed
	}
	newOrder.Items = items
	errSaving := oc.orders.Save(ctx, newOrder)
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
