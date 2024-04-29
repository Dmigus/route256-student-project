// Package orders содержит реализацию заказов только для чтения из PostgreSQL.
package orders

import (
	"context"
	"errors"
	"route256.ozon.ru/project/loms/internal/pkg/sqlmetrics"

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
		RecordDuration(table string, category sqlmetrics.SQLCategory, f func() error)
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

// LoadAll загружает информацию о всех заказах из БД PostgreSQL
func (po *Orders) LoadAll(ctx context.Context) ([]*models.Order, error) {
	orders, err := po.loadAllOrdersRowWithoutItems(ctx)
	if err != nil {
		return nil, err
	}
	items, err := po.loadItemsForAllOrders(ctx)
	if err != nil {
		return nil, err
	}
	for _, order := range orders {
		order.Items = items[order.Id()]
	}
	return orders, nil
}

// loadAllOrdersRowWithoutItems загружает информацию о заказах без их товаров
func (po *Orders) loadAllOrdersRowWithoutItems(ctx context.Context) ([]*models.Order, error) {
	var orders []Order
	var err error
	po.reqDur.RecordDuration(orderTableName, sqlmetrics.Select, func() error {
		orders, err = po.queries.selectAllOrders(ctx)
		return err
	})
	if err != nil {
		return nil, err
	}
	modelOrders := make([]*models.Order, 0, len(orders))
	for _, order := range orders {
		modelOrder := models.NewOrder(order.UserID, order.ID)
		modelOrder.IsItemsReserved = order.AreItemsReserved
		modelOrder.Status = orderStatusFromString(order.Status)
		modelOrders = append(modelOrders, modelOrder)
	}
	return modelOrders, nil
}

// loadItemsForAllOrders загружает информацию о товарах заказов
func (po *Orders) loadItemsForAllOrders(ctx context.Context) (map[int64][]models.OrderItem, error) {
	var items []OrderItem
	var err error
	po.reqDur.RecordDuration(orderItemTableName, sqlmetrics.Select, func() error {
		items, err = po.queries.selectAllOrdersItems(ctx)
		return err
	})
	if err != nil {
		return nil, err
	}
	orderNumWithItems := make(map[int64][]models.OrderItem)
	for _, it := range items {
		item := models.OrderItem{SkuId: it.SkuID, Count: uint16(it.Count)}
		orderNumWithItems[it.OrderID] = append(orderNumWithItems[it.OrderID], item)
	}
	return orderNumWithItems, nil
}

func (po *Orders) loadOrderRowWithoutItems(ctx context.Context, orderID int64) (*models.Order, error) {
	var row selectOrderRow
	var err error
	po.reqDur.RecordDuration(orderTableName, sqlmetrics.Select, func() error {
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
	po.reqDur.RecordDuration(orderItemTableName, sqlmetrics.Select, func() error {
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
