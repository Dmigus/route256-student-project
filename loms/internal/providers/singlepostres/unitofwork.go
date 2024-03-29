package singlepostres

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type UnitOfWork struct {
	conn txBeginner2
}

func NewUnitOfWork(conn txBeginner2) *UnitOfWork {
	return &UnitOfWork{conn: conn}
}

type txBeginner2 interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

func (u *UnitOfWork) Transactional(ctx context.Context, f func(ctx context.Context, orders any, stocks any) error) error {
	tx, err := u.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	orders := NewPostgresOrders2(tx)
	stocks := NewPostgresStocks2(tx)
	err = f(ctx, orders, stocks)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
