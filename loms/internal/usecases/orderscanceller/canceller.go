package orderscanceller

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"route256.ozon.ru/project/loms/internal/models"
)

var errWrongOrderStatus = errors.Wrap(models.ErrFailedPrecondition, "order status is wrong")

type orderRepo interface {
	Save(context.Context, *models.Order) error
	Load(context.Context, int64) (*models.Order, error)
}

type stockCanceller interface {
	CancelReserved(context.Context, []models.OrderItem) error
	AddItems(context.Context, []models.OrderItem) error
}

type OrderCanceller struct {
	orders orderRepo
	stocks stockCanceller
}

func NewOrderCanceller(orders orderRepo, stocks stockCanceller) *OrderCanceller {
	return &OrderCanceller{orders: orders, stocks: stocks}
}

func (oc *OrderCanceller) Cancel(ctx context.Context, orderId int64) error {
	order, err := oc.orders.Load(ctx, orderId)
	if err != nil {
		return fmt.Errorf("could not load order %d: %w", orderId, err)
	}
	if order.Status == models.Cancelled {
		return errWrongOrderStatus
	}
	if order.IsItemsReserved {
		if err := oc.cancelReserved(ctx, order); err != nil {
			return err
		}
	} else if order.Status == models.Payed {
		if err := oc.cancelPayed(ctx, order); err != nil {
			return err
		}
	} else {
		order.Status = models.Cancelled
	}
	if err = oc.orders.Save(ctx, order); err != nil {
		return fmt.Errorf("could not save order %d: %w", orderId, err)
	}
	return nil
}

func (oc *OrderCanceller) cancelReserved(ctx context.Context, order *models.Order) error {
	if order.IsItemsReserved {
		err := oc.stocks.CancelReserved(ctx, order.Items)
		if err != nil {
			return fmt.Errorf("could not cancel reserved items for order %d: %w", order.Id(), err)
		}
		order.IsItemsReserved = false
		order.Status = models.Cancelled
	}
	return nil
}

func (oc *OrderCanceller) cancelPayed(ctx context.Context, order *models.Order) error {
	if order.Status == models.Payed {
		err := oc.stocks.AddItems(ctx, order.Items)
		if err != nil {
			return fmt.Errorf("could not return payed items for order %d: %w", order.Id(), err)
		}
		order.Status = models.Cancelled
	}
	return nil
}
