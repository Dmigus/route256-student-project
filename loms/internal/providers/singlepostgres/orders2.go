package singlepostgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"route256.ozon.ru/project/loms/internal/models"
)

type connect interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
}

type PostgresOrders2 struct {
	tx connect
}

func NewPostgresOrders2(tx connect) *PostgresOrders2 {
	return &PostgresOrders2{tx: tx}
}

// Create создаёт заказ для юзера userID и товарами items в репозитории и возращает его
func (po *PostgresOrders2) Create(ctx context.Context, userID int64, items []models.OrderItem) (*models.Order, error) {
	var orderId int64
	row := po.tx.QueryRow(ctx, createOrder, userID, orderStatusToString(models.New), false)
	if err := row.Scan(&orderId); err != nil {
		return nil, err
	}
	order := models.NewOrder(userID, orderId)
	for _, it := range items {
		_, err := po.tx.Exec(ctx, insertOrderItem, orderId, it.SkuId, it.Count)
		if err != nil {
			return nil, err
		}
	}
	order.Items = items
	return order, nil
}

// Save сохраняет заказ в БД в PostgreSQL. Если его не было, то создаётся новый. Если был, то обновляется. Изменение позиций заказа после создания не предусмотрено
func (po *PostgresOrders2) Save(ctx context.Context, order *models.Order) error {
	_, err := po.tx.Exec(ctx, updateOrder, order.Id(), orderStatusToString(order.Status), order.IsItemsReserved)
	return err
}

// Load загружает информацию о заказе из БД в PostgreSQL
func (po *PostgresOrders2) Load(ctx context.Context, orderID int64) (*models.Order, error) {
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

func (po *PostgresOrders2) loadOrderRowWithoutItems(ctx context.Context, orderID int64) (*models.Order, error) {
	var userID int64
	var strStatus string
	var isItemsReserved bool
	row := po.tx.QueryRow(ctx, selectOrder, orderID)
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

func (po *PostgresOrders2) readItemsForOrder(ctx context.Context, orderID int64) ([]models.OrderItem, error) {
	rows, err := po.tx.Query(ctx, selectOrderItems, orderID)
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
