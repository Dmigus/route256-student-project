package singlepostres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"route256.ozon.ru/project/loms/internal/models"
)

// Предполагается, что транзакционная логика каждого юзкейса примитивна: всё в одной транзакции. Если нет ошибки - коммитим. Если есть ошибка полностью откатываем.
type usecasesToWrap interface {
	CreateOrder(ctx context.Context, userId int64, items []models.OrderItem) (int64, error)
	PayOrder(ctx context.Context, orderId int64) error
	GetNumOfAvailable(ctx context.Context, skuId int64) (uint64, error)
	GetOrder(ctx context.Context, orderId int64) (*models.Order, error)
	CancelOrder(ctx context.Context, orderId int64) error
}

var _ usecasesToWrap = (*TrWrapper)(nil)

type TrWrapper struct {
	usecases usecasesToWrap
	pool     *pgxpool.Pool
}

func (t *TrWrapper) InTx(ctx context.Context, txOptions pgx.TxOptions, f func(ctx context.Context) error) error {
	return InTx(t.pool, ctx, txOptions, f)
}

func (t *TrWrapper) CreateOrder(ctx context.Context, userId int64, items []models.OrderItem) (orderId int64, err error) {
	err = InTx(t.pool, ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(ctx context.Context) error {
		orderId, err = t.usecases.CreateOrder(ctx, userId, items)
		return err
	})
	if err != nil {
		return 0, err
	}
	return orderId, nil

	//tx, err := t.pool.Begin(ctx)
	//if err != nil {
	//	return 0, err
	//}
	//defer tx.Rollback(ctx)
	//trCtx := context.WithValue(ctx, trKey, tx)
	//orderId, err = t.usecases.CreateOrder(trCtx, userId, items)
	//if err != nil {
	//	return 0, err
	//}
	//err = tx.Commit(ctx)
	//if err != nil {
	//	return 0, err
	//}
	//return orderId, nil
}

func (t *TrWrapper) PayOrder(ctx context.Context, orderId int64) error {
	return InTx(t.pool, ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead}, func(ctx context.Context) error {
		return t.usecases.PayOrder(ctx, orderId)
	})
}

func (t *TrWrapper) GetNumOfAvailable(ctx context.Context, skuId int64) (cnt uint64, err error) {
	err = InTx(t.pool, ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(ctx context.Context) error {
		cnt, err = t.usecases.GetNumOfAvailable(ctx, skuId)
		return err
	})
	if err != nil {
		return 0, err
	}
	return cnt, nil

	//tx, err := t.pool.Begin(ctx)
	//if err != nil {
	//	return 0, err
	//}
	//defer tx.Rollback(ctx)
	//trCtx := context.WithValue(ctx, trKey, tx)
	//cnt, err := t.usecases.GetNumOfAvailable(trCtx, skuId)
	//if err != nil {
	//	return 0, err
	//}
	//err = tx.Commit(ctx)
	//if err != nil {
	//	return 0, err
	//}
	//return cnt, nil
}

func (t *TrWrapper) GetOrder(ctx context.Context, orderId int64) (order *models.Order, err error) {
	err = InTx(t.pool, ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(ctx context.Context) error {
		order, err = t.usecases.GetOrder(ctx, orderId)
		return err
	})
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (t *TrWrapper) CancelOrder(ctx context.Context, orderId int64) (err error) {
	return InTx(t.pool, ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead}, func(ctx context.Context) error {
		return t.usecases.CancelOrder(ctx, orderId)
	})
}

func NewTrWrapper(usecases usecasesToWrap, pool *pgxpool.Pool) *TrWrapper {
	return &TrWrapper{usecases: usecases, pool: pool}
}
