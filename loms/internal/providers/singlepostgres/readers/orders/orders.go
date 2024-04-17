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

const (
	orderTableName     = "order"
	orderItemTableName = "order_item"
)

// Orders представялет реализацию репозитория заказов с методами для чтения данных
type (
	durationRecorder interface {
		RecordDuration(table, category string, f func() error)
	}
	Orders struct {
		queries *Queries
		reqDur  durationRecorder
	}
)

// NewOrders создаёт объект репозитория заказов, работающего в рамках транзакции tx
func NewOrders(tx DBTX, reqDur durationRecorder) *Orders {
	return &Orders{queries: New(tx), reqDur: reqDur}
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
	var row selectOrderRow
	var err error
	po.reqDur.RecordDuration(orderTableName, "select", func() error {
		row, err = po.queries.selectOrder(ctx, orderID)
		return err
	})
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
	var rows []selectOrderItemsRow
	var err error
	po.reqDur.RecordDuration(orderItemTableName, "select", func() error {
		rows, err = po.queries.selectOrderItems(ctx, orderID)
		return err
	})
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
