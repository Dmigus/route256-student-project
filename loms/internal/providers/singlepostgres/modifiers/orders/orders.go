// Package orders содержит реализацию заказов для транзакционной модификации данных в PostgreSQL.
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

// Orders представялет реализацию репозитория заказов с методами для модификации данных
type (
	durationRecorder interface {
		RecordDuration(table string, category sqlmetrics.SQLCategory, f func() error)
	}
	Orders struct {
		queries *Queries
		durRec  durationRecorder
	}
)

// NewOrders создаёт объект репозитория заказов, работающего в рамках транзакции tx
func NewOrders(tx DBTX, durRec durationRecorder) *Orders {
	return &Orders{queries: New(tx), durRec: durRec}
}

// Create создаёт заказ для юзера userID и товарами items в репозитории и возращает его
func (po *Orders) Create(ctx context.Context, userID int64, items []models.OrderItem) (*models.Order, error) {
	params := createOrderParams{UserID: userID, Status: orderStatusToString(models.New), AreItemsReserved: false}
	var orderID int64
	var err error
	po.durRec.RecordDuration(orderTableName, sqlmetrics.Insert, func() error {
		orderID, err = po.queries.createOrder(ctx, params)
		return err
	})
	if err != nil {
		return nil, err
	}
	order := models.NewOrder(userID, orderID)
	itemsParams := insertItemParamsFrom(orderID, items)
	po.durRec.RecordDuration(orderItemTableName, sqlmetrics.Insert, func() error {
		_, err = po.queries.insertOrderItem(ctx, itemsParams)
		return err
	})
	if err != nil {
		return nil, err
	}
	order.Items = items
	return order, nil
}

func insertItemParamsFrom(orderID int64, items []models.OrderItem) []insertOrderItemParams {
	itemsParams := make([]insertOrderItemParams, 0, len(items))
	for _, it := range items {
		params := insertOrderItemParams{OrderID: orderID, SkuID: it.SkuId, Count: int32(it.Count)}
		itemsParams = append(itemsParams, params)
	}
	return itemsParams
}

// Save сохраняет заказ в БД в PostgreSQL. Изменение позиций заказа не предусмотрено
func (po *Orders) Save(ctx context.Context, order *models.Order) error {
	params := updateOrderParams{
		ID:               order.Id(),
		Status:           orderStatusToString(order.Status),
		AreItemsReserved: order.IsItemsReserved,
	}
	var err error
	po.durRec.RecordDuration(orderTableName, sqlmetrics.Update, func() error {
		err = po.queries.updateOrder(ctx, params)
		return err
	})
	return err
}

func orderStatusToString(os models.OrderStatus) string {
	switch os {
	case models.New:
		return "New"
	case models.AwaitingPayment:
		return "AwaitingPayment"
	case models.Failed:
		return "Failed"
	case models.Payed:
		return "Payed"
	case models.Cancelled:
		return "Cancelled"
	default:
		return "Undefined"
	}
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

// Load загружает информацию о заказе из БД в PostgreSQL, производя SELECT FOR UPDATE
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
	po.durRec.RecordDuration(orderTableName, sqlmetrics.Select, func() error {
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
	po.durRec.RecordDuration(orderItemTableName, sqlmetrics.Select, func() error {
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
