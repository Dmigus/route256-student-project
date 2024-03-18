package orderscanceller

import (
	"context"
	"fmt"
	"route256.ozon.ru/project/loms/internal/models"
)

var errWrongOrderStatus = fmt.Errorf("order status is not awaiting")

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

func (oc *OrderCanceller) Cancel(ctx context.Context, orderId int64) error {
	order, err := oc.orders.Load(ctx, orderId)
	if err != nil {
		return err
	}
	if order.Status != models.AwaitingPayment {
		return errWrongOrderStatus
	}
	err = oc.stocks.CancelReserved(ctx, order.Items)
	if err != nil {
		return err
	}
	order.Status = models.Cancelled
	if err = oc.orders.Save(ctx, order); err != nil {
		return err
	}
	return nil
}
