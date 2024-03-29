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
	uow unitOfWork
}

type unitOfWork interface {
	Transactional(context.Context, func(ctx context.Context, orders any, stocks any) error) error
}

func NewOrderCanceller(uow unitOfWork) *OrderCanceller {
	return &OrderCanceller{uow: uow}
}

func (oc *OrderCanceller) Cancel(ctx context.Context, orderId int64) error {
	return oc.uow.Transactional(ctx, func(ctx context.Context, ordersAny any, stocksAny any) error {
		orders := ordersAny.(orderRepo)
		stocks := stocksAny.(stockCanceller)
		order, err := orders.Load(ctx, orderId)
		if err != nil {
			return fmt.Errorf("could not load order %d: %w", orderId, err)
		}
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
		if err = orders.Save(ctx, order); err != nil {
			return fmt.Errorf("could not save order %d: %w", orderId, err)
		}
		return nil
	})
}

func cancelReserved(ctx context.Context, stocks stockCanceller, order *models.Order) error {
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

func cancelPayed(ctx context.Context, stocks stockCanceller, order *models.Order) error {
	if order.Status == models.Payed {
		err := stocks.AddItems(ctx, order.Items)
		if err != nil {
			return fmt.Errorf("could not return payed items for order %d: %w", order.Id(), err)
		}
		order.Status = models.Cancelled
	}
	return nil
}
