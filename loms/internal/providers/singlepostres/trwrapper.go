package singlepostres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"route256.ozon.ru/project/loms/internal/models"
)

// Предполагается, что транзакционная логика каждого юзкейса примитивна: всё в одной транзакции. Если нет ошибки - коммитим. Если есть ошибка полностью откатываем.
type usecasesToWrap interface {
	CreateOrder(ctx context.Context, userID int64, items []models.OrderItem) (int64, error)
	PayOrder(ctx context.Context, orderID int64) error
	GetNumOfAvailable(ctx context.Context, skuID int64) (uint64, error)
	GetOrder(ctx context.Context, orderID int64) (*models.Order, error)
	CancelOrder(ctx context.Context, orderID int64) error
}

var _ usecasesToWrap = (*TrWrapper)(nil)

// TrWrapper это обёртка, которая вызывает юзкейсы в транзакциях
type TrWrapper struct {
	usecases usecasesToWrap
	conn     txBeginner
}

// NewTrWrapper создаёт экземпляр обёртки над юзкейсами
func NewTrWrapper(usecases usecasesToWrap, conn txBeginner) *TrWrapper {
	return &TrWrapper{usecases: usecases, conn: conn}
}

// CreateOrder исполняет usecases.CreateOrder в транзакции с уровнем ReadCommitted
func (t *TrWrapper) CreateOrder(ctx context.Context, userID int64, items []models.OrderItem) (orderID int64, err error) {
	err = InTx(ctx, t.conn, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(ctx context.Context) error {
		orderID, err = t.usecases.CreateOrder(ctx, userID, items)
		return err
	})
	if err != nil {
		return 0, err
	}
	return orderID, nil
}

// PayOrder исполняет usecases.PayOrder в транзакции с уровнем RepeatableRead
func (t *TrWrapper) PayOrder(ctx context.Context, orderID int64) error {
	return InTx(ctx, t.conn, pgx.TxOptions{IsoLevel: pgx.RepeatableRead}, func(ctx context.Context) error {
		return t.usecases.PayOrder(ctx, orderID)
	})
}

// GetNumOfAvailable исполняет usecases.GetNumOfAvailable в транзакции с уровнем ReadCommitted
func (t *TrWrapper) GetNumOfAvailable(ctx context.Context, skuID int64) (cnt uint64, err error) {
	err = InTx(ctx, t.conn, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(ctx context.Context) error {
		cnt, err = t.usecases.GetNumOfAvailable(ctx, skuID)
		return err
	})
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

// GetOrder исполняет usecases.GetOrder в транзакции с уровнем ReadCommitted
func (t *TrWrapper) GetOrder(ctx context.Context, orderID int64) (order *models.Order, err error) {
	err = InTx(ctx, t.conn, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(ctx context.Context) error {
		order, err = t.usecases.GetOrder(ctx, orderID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return order, nil
}

// CancelOrder исполняет usecases.CancelOrder в транзакции с уровнем RepeatableRead
func (t *TrWrapper) CancelOrder(ctx context.Context, orderID int64) (err error) {
	return InTx(ctx, t.conn, pgx.TxOptions{IsoLevel: pgx.RepeatableRead}, func(ctx context.Context) error {
		return t.usecases.CancelOrder(ctx, orderID)
	})
}
