package orderspayer

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"route256.ozon.ru/project/loms/internal/models"
)

var errWrongOrderStatus = errors.Wrap(models.ErrFailedPrecondition, "order status is wrong")

type (
	OrderRepo interface {
		Save(context.Context, *models.Order) error
		Load(context.Context, int64) (*models.Order, error)
	}
	StockRepo interface {
		RemoveReserved(context.Context, []models.OrderItem) error
	}
	txManager interface {
		WithinTransaction(context.Context, func(ctx context.Context, orders OrderRepo, stocks StockRepo) error) error
	}
	OrdersPayer struct {
		tx txManager
	}
)

func NewOrdersPayer(tx txManager) *OrdersPayer {
	return &OrdersPayer{tx: tx}
}

func (or *OrdersPayer) Pay(ctx context.Context, orderId int64) error {
	return or.tx.WithinTransaction(ctx, func(ctx context.Context, orders OrderRepo, stocks StockRepo) error {
		order, err := orders.Load(ctx, orderId)
		if err != nil {
			return fmt.Errorf("could not load order %d: %w", orderId, err)
		}
		if order.Status != models.AwaitingPayment {
			return errWrongOrderStatus
		}
		err = stocks.RemoveReserved(ctx, order.Items)
		if err != nil {
			return fmt.Errorf("could not remove reserved items for order %d: %w", orderId, err)
		}
		order.Status = models.Payed
		order.IsItemsReserved = false
		if err = orders.Save(ctx, order); err != nil {
			return fmt.Errorf("could not save order %d: %w", orderId, err)
		}
		return nil
	})
}
