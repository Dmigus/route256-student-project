package singlepostres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"route256.ozon.ru/project/loms/internal/models"
)

var errOrderNotFound = errors.Wrap(models.ErrNotFound, "order is not found")

type PostgresOrders struct {
}

const (
	updateOrder      = `UPDATE "order" SET status = $2, are_items_reserved = $3 where id = $1`
	createOrder      = `INSERT INTO "order"(id, user_id, status, are_items_reserved) VALUES ($1, $2, $3, $4)`
	insertOrderItem  = `INSERT INTO order_item(order_id, sku_id, count) VALUES ($1, $2, $3)`
	selectOrder      = `SELECT user_id, status, are_items_reserved from "order" where id = $1`
	selectOrderItems = `SELECT sku_id, count FROM order_item WHERE order_id = $1`
)

// Save сохраняет заказ. Если его не было, то создаётся новый. Если был, то обновляется. Изменение позиций заказа после создания не предусмотрено
func (po *PostgresOrders) Save(ctx context.Context, order *models.Order) error {
	tx := ctx.Value(trKey).(pgx.Tx)
	tag, err := tx.Exec(ctx, updateOrder, order.Id(), orderStatusToString(order.Status), order.IsItemsReserved)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return po.createNewOrder(ctx, order)
	}
	return nil
}

func (po *PostgresOrders) createNewOrder(ctx context.Context, order *models.Order) error {
	tx := ctx.Value(trKey).(pgx.Tx)
	_, err := tx.Exec(ctx, createOrder, order.Id(), order.UserId, orderStatusToString(order.Status), order.IsItemsReserved)
	if err != nil {
		return err
	}
	for _, it := range order.Items {
		_, err = tx.Exec(ctx, insertOrderItem, order.Id(), it.SkuId, it.Count)
		if err != nil {
			return err
		}
	}
	return nil
}

func (po *PostgresOrders) Load(ctx context.Context, orderId int64) (*models.Order, error) {
	order, err := loadOrderRowWithoutItems(ctx, orderId)
	if err != nil {
		return nil, err
	}
	orderItems, err := readItemsForOrder(ctx, orderId)
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

func loadOrderRowWithoutItems(ctx context.Context, orderId int64) (*models.Order, error) {
	var userId int64
	var strStatus string
	var isItemsReserved bool
	tx := ctx.Value(trKey).(pgx.Tx)
	row := tx.QueryRow(ctx, selectOrder, orderId)
	err := row.Scan(&userId, &strStatus, &isItemsReserved)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = errOrderNotFound
		}
		return nil, err
	}
	order := models.NewOrder(userId, orderId)
	order.Status = orderStatusFromString(strStatus)
	order.IsItemsReserved = isItemsReserved
	return order, nil
}

func readItemsForOrder(ctx context.Context, orderId int64) ([]models.OrderItem, error) {
	tx := ctx.Value(trKey).(pgx.Tx)
	rows, err := tx.Query(ctx, selectOrderItems, orderId)
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
	var skuId int64
	var cnt uint16
	err := rows.Scan(&skuId, &cnt)
	if err != nil {
		return models.OrderItem{}, err
	}
	return models.OrderItem{SkuId: skuId, Count: cnt}, nil
}
