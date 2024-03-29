package read

import (
	"context"
	"route256.ozon.ru/project/loms/internal/models"
)

type Orders struct {
	queries *Queries
}

func NewOrders(tx DBTX) *Orders {
	return &Orders{queries: New(tx)}
}

func orderStatusFromDTO(s OrderStatus) models.OrderStatus {
	switch s {
	case OrderStatusNew:
		return models.New
	case OrderStatusAwaitingPayment:
		return models.AwaitingPayment
	case OrderStatusFailed:
		return models.Failed
	case OrderStatusPayed:
		return models.Payed
	case OrderStatusCancelled:
		return models.Cancelled
	default:
		return models.OrderStatus(0)
	}
}

// Load загружает информацию о заказе из БД в PostgreSQL
func (po *Orders) Load(ctx context.Context, orderID int64) (*models.Order, error) {
	order, err := po.loadOrderRowWithoutItems(ctx, orderID)
	if err != nil {
		return nil, err
	}
	orderItems, err := po.readItemsForOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}
	order.Items = orderItems
	return order, nil
}

func (po *Orders) loadOrderRowWithoutItems(ctx context.Context, orderID int64) (*models.Order, error) {
	row, err := po.queries.selectOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}
	order := models.NewOrder(row.UserID, orderID)
	order.Status = orderStatusFromDTO(row.Status)
	order.IsItemsReserved = row.AreItemsReserved
	return order, nil
}

func (po *Orders) readItemsForOrder(ctx context.Context, orderID int64) ([]models.OrderItem, error) {
	rows, err := po.queries.selectOrderItems(ctx, orderID)
	if err != nil {
		return nil, err
	}
	items := make([]models.OrderItem, 0, len(rows))
	for _, it := range rows {
		item := models.OrderItem{SkuId: it.SkuID, Count: uint16(it.Count)}
		items = append(items, item)
	}
	return items, nil
}
