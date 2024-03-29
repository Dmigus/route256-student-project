package singlepostres

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type UnitOfWork[O, S any] struct {
	conn txBeginner2
}

func NewUnitOfWork[O, S any](conn txBeginner2) *UnitOfWork[O, S] {
	return &UnitOfWork[O, S]{conn: conn}
}

type txBeginner2 interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

func (u *UnitOfWork[O, S]) Transactional(ctx context.Context, f func(ctx context.Context, orders O, stocks S) error) error {
	tx, err := u.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var orders any = NewPostgresOrders2(tx)
	var stocks any = NewPostgresStocks2(tx)
	err = f(ctx, orders.(O), stocks.(S))
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
