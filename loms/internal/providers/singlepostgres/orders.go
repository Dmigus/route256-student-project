package singlepostgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"route256.ozon.ru/project/loms/internal/models"
)

var errOrderNotFound = errors.Wrap(models.ErrNotFound, "order is not found")

// PostgresOrders это реализация репозитория заказов для использования с БД в PostgreSQL
type PostgresOrders struct {
}

const (
	updateOrder      = `UPDATE "order" SET status = $2, are_items_reserved = $3 where id = $1`
	createOrder      = `INSERT INTO "order"(user_id, status, are_items_reserved) VALUES ($1, $2, $3) RETURNING id`
	insertOrderItem  = `INSERT INTO order_item(order_id, sku_id, count) VALUES ($1, $2, $3)`
	selectOrder      = `SELECT user_id, status, are_items_reserved from "order" where id = $1`
	selectOrderItems = `SELECT sku_id, count FROM order_item WHERE order_id = $1`
)

// Create создаёт заказ для юзера userID и товарами items в репозитории и возращает его
func (po *PostgresOrders) Create(ctx context.Context, userID int64, items []models.OrderItem) (*models.Order, error) {
	tx := ctx.Value(trKey).(pgx.Tx)
	var orderId int64
	row := tx.QueryRow(ctx, createOrder, userID, orderStatusToString(models.New), false)
	if err := row.Scan(&orderId); err != nil {
		return nil, err
	}
	order := models.NewOrder(userID, orderId)
	for _, it := range items {
		_, err := tx.Exec(ctx, insertOrderItem, orderId, it.SkuId, it.Count)
		if err != nil {
			return nil, err
		}
	}
	order.Items = items
	return order, nil
}

// Save сохраняет заказ в БД в PostgreSQL. Если его не было, то создаётся новый. Если был, то обновляется. Изменение позиций заказа после создания не предусмотрено
func (po *PostgresOrders) Save(ctx context.Context, order *models.Order) error {
	tx := ctx.Value(trKey).(pgx.Tx)
	_, err := tx.Exec(ctx, updateOrder, order.Id(), orderStatusToString(order.Status), order.IsItemsReserved)
	return err
}

// Load загружает информацию о заказе из БД в PostgreSQL
func (po *PostgresOrders) Load(ctx context.Context, orderID int64) (*models.Order, error) {
	order, err := loadOrderRowWithoutItems(ctx, orderID)
	if err != nil {
		return nil, err
	}
	orderItems, err := readItemsForOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}
	order.Items = orderItems
	return order, nil
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

func loadOrderRowWithoutItems(ctx context.Context, orderID int64) (*models.Order, error) {
	var userID int64
	var strStatus string
	var isItemsReserved bool
	tx := ctx.Value(trKey).(pgx.Tx)
	row := tx.QueryRow(ctx, selectOrder, orderID)
	err := row.Scan(&userID, &strStatus, &isItemsReserved)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = errOrderNotFound
		}
		return nil, err
	}
	order := models.NewOrder(userID, orderID)
	order.Status = orderStatusFromString(strStatus)
	order.IsItemsReserved = isItemsReserved
	return order, nil
}

func readItemsForOrder(ctx context.Context, orderID int64) ([]models.OrderItem, error) {
	tx := ctx.Value(trKey).(pgx.Tx)
	rows, err := tx.Query(ctx, selectOrderItems, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orderItems []models.OrderItem
	for rows.Next() {
		item, err := readNextOrderItem(rows)
		if err != nil {
			return nil, err
		}
		orderItems = append(orderItems, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orderItems, nil
}

func readNextOrderItem(rows pgx.Rows) (models.OrderItem, error) {
	var skuID int64
	var cnt uint16
	err := rows.Scan(&skuID, &cnt)
	if err != nil {
		return models.OrderItem{}, err
	}
	return models.OrderItem{SkuId: skuID, Count: cnt}, nil
}
