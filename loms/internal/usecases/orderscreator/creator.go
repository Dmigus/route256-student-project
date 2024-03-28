package orderscreator

import (
	"context"
	"errors"
	"fmt"
	"route256.ozon.ru/project/loms/internal/models"
)

type stocksStorage interface {
	Reserve(context.Context, []models.OrderItem) error
}

type ordersStorage interface {
	Create(context.Context, int64, []models.OrderItem) (*models.Order, error)
	Save(context.Context, *models.Order) error
}

type OrdersCreator struct {
	orders ordersStorage
	stocks stocksStorage
}

func NewOrdersCreator(orders ordersStorage, stocks stocksStorage) *OrdersCreator {
	return &OrdersCreator{orders: orders, stocks: stocks}
}

func (oc *OrdersCreator) Create(ctx context.Context, userID int64, items []models.OrderItem) (int64, error) {
	order, err := oc.orders.Create(ctx, userID, items)
	if err != nil {
		return 0, fmt.Errorf("could not create new order for user %d: %w", userID, err)
	}
	errReserving := oc.stocks.Reserve(ctx, items)
	if errReserving != nil {
		errReserving = fmt.Errorf("could not reserve items for user %d: %w", userID, errReserving)
		order.Status = models.Failed
	} else {
		order.Status = models.AwaitingPayment
		order.IsItemsReserved = true
	}
	order.Items = items
	errSaving := oc.orders.Save(ctx, order)
	if errSaving != nil {
		errSaving = fmt.Errorf("could not save created order for user %d: %w", userID, errSaving)
	}
	errs := errors.Join(errSaving, errReserving)
	if errs != nil {
		return 0, errs
	}
	return order.Id(), nil
}
