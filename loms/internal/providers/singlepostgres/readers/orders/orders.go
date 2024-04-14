// Package orders содержит реализацию заказов только для чтения из PostgreSQL.
package orders

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	pkgerrors "github.com/pkg/errors"

	"route256.ozon.ru/project/loms/internal/models"
)

var errOrderNotFound = pkgerrors.Wrap(models.ErrNotFound, "order is not found")

// Orders представялет реализацию репозитория заказов с методами для чтения данных
type Orders struct {
	queries *Queries
}

// NewOrders создаёт объект репозитория заказов, работающего в рамках транзакции tx
func NewOrders(tx DBTX) *Orders {
	return &Orders{queries: New(tx)}
}

func orderStatusFromString(s string) models.OrderStatus {
	switch s {
	case "New":
		return models.New
	case "AwaitingPayment":
		return models.AwaitingPayment
	case "Failed":
		return models.Failed
	case "Payed":
		return models.Payed
	case "Cancelled":
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
		if errors.Is(err, pgx.ErrNoRows) {
			err = errOrderNotFound
		}
		return nil, err
	}
	order := models.NewOrder(row.UserID, orderID)
	order.Status = orderStatusFromString(row.Status)
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
