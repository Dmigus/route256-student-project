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

func (oc *OrdersCreator) Create(ctx context.Context, _ int64, items []models.OrderItem) (int64, error) {
	newOrder := oc.createOrderInstance()
	errReserving := oc.stocks.Reserve(ctx, items)
	if errReserving == nil {
		newOrder.Status = models.AwaitingPayment
	} else {
		newOrder.Status = models.Failed
	}
	errSaving := oc.orders.Save(ctx, newOrder)
	errs := errors.Join(errSaving, errReserving)
	if errs != nil {
		return 0, errs
	}
	return newOrder.Id(), nil
}

func (oc *OrdersCreator) createOrderInstance() *models.Order {
	newOrderId := oc.orderIdGenerator.NewId()
	return models.NewOrder(newOrderId)
}
