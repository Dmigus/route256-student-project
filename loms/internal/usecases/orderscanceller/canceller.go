package orderscanceller

import (
	"context"
	"route256.ozon.ru/project/loms/internal/models"
)

type orderRepo interface {
	Save(context.Context, *models.Order) error
	Load(context.Context, int64) (*models.Order, error)
}

type stockCanceller interface {
	CancelReserved(context.Context, []models.OrderItem) error
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
		return err
	}
	if order.Status == models.Cancelled {
		return models.ErrWrongOrderStatus
	}
	if order.IsItemsReserved {
		err = oc.stocks.CancelReserved(ctx, order.Items)
		if err != nil {
			return err
		}
		order.IsItemsReserved = false
	}
	order.Status = models.Cancelled
	if err = oc.orders.Save(ctx, order); err != nil {
		return err
	}
	return nil
}
