package orderspayer

import (
	"context"
	"route256.ozon.ru/project/loms/internal/models"
)

type orderRepo interface {
	Save(context.Context, *models.Order) error
	Load(context.Context, int64) (*models.Order, error)
}

type stockRemover interface {
	RemoveReserved(context.Context, []models.OrderItem) error
}

type OrdersPayer struct {
	orders orderRepo
	stocks stockRemover
}

func NewOrdersPayer(orders orderRepo, stocks stockRemover) *OrdersPayer {
	return &OrdersPayer{orders: orders, stocks: stocks}
}

func (or *OrdersPayer) Pay(ctx context.Context, orderId int64) error {
	order, err := or.orders.Load(ctx, orderId)
	if err != nil {
		return err
	}
	if order.Status != models.AwaitingPayment {
		return models.ErrWrongOrderStatus
	}
	err = or.stocks.RemoveReserved(ctx, order.Items)
	if err != nil {
		return err
	}
	order.Status = models.Payed
	order.IsItemsReserved = false
	if err = or.orders.Save(ctx, order); err != nil {
		return err
	}
	return nil
}
