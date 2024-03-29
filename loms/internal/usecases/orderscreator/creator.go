package orderscreator

import (
	"context"
	"errors"
	"fmt"
	"route256.ozon.ru/project/loms/internal/models"
)

type (
	StockRepo interface {
		Reserve(context.Context, []models.OrderItem) error
	}

	OrderRepo interface {
		Create(context.Context, int64, []models.OrderItem) (*models.Order, error)
		Save(context.Context, *models.Order) error
	}
	txManager interface {
		WithinTransaction(context.Context, func(ctx context.Context, orders OrderRepo, stocks StockRepo) error) error
	}

	OrdersCreator struct {
		tx txManager
	}
)

func NewOrdersCreator(tx txManager) *OrdersCreator {
	return &OrdersCreator{tx: tx}
}

func (oc *OrdersCreator) Create(ctx context.Context, userID int64, items []models.OrderItem) (orderID int64, err error) {
	err = oc.tx.WithinTransaction(ctx, func(ctx context.Context, orders OrderRepo, stocks StockRepo) error {
		orderID, err = createOrder(ctx, userID, items, orders, stocks)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return orderID, nil
}

func createOrder(ctx context.Context, userID int64, items []models.OrderItem, orders OrderRepo, stocks StockRepo) (int64, error) {
	order, err := orders.Create(ctx, userID, items)
	if err != nil {
		return 0, fmt.Errorf("could not create new order for user %d: %w", userID, err)
	}
	errReserving := stocks.Reserve(ctx, items)
	if errReserving != nil {
		errReserving = fmt.Errorf("could not reserve items for user %d: %w", userID, errReserving)
		order.Status = models.Failed
	} else {
		order.Status = models.AwaitingPayment
		order.IsItemsReserved = true
	}
	errSaving := orders.Save(ctx, order)
	if errSaving != nil {
		errSaving = fmt.Errorf("could not save created order for user %d: %w", userID, errSaving)
	}
	errs := errors.Join(errSaving, errReserving)
	if errs != nil {
		return 0, errs
	}
	return order.Id(), nil
}
