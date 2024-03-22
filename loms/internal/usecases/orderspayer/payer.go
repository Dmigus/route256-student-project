package orderspayer

import (
	"context"
	"fmt"
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
		return fmt.Errorf("could not load order %d: %w", orderId, err)
	}
	if order.Status != models.AwaitingPayment {
		return models.ErrWrongOrderStatus
	}
	err = or.stocks.RemoveReserved(ctx, order.Items)
	if err != nil {
		return fmt.Errorf("could not remove reserved items for order %d: %w", orderId, err)
	}
	order.Status = models.Payed
	order.IsItemsReserved = false
	if err = or.orders.Save(ctx, order); err != nil {
		return fmt.Errorf("could not save order %d: %w", orderId, err)
	}
	return nil
}
