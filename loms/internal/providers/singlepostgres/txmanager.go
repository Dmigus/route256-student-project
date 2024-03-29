package singlepostgres

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type (
	TxBeginner interface {
		Begin(ctx context.Context) (pgx.Tx, error)
	}
	TxManager[OrderRepoType, StocksRepoType any] struct {
		conn          TxBeginner
		ordersCreator func(conn TxBeginner) OrderRepoType
		stocksCreator func(conn TxBeginner) StocksRepoType
	}
)

func NewTxManager[OrderRepoType, StocksRepoType any](conn TxBeginner, ordersCreator func(conn TxBeginner) OrderRepoType, stocksCreator func(conn TxBeginner) StocksRepoType) *TxManager[OrderRepoType, StocksRepoType] {
	return &TxManager[OrderRepoType, StocksRepoType]{
		conn:          conn,
		ordersCreator: ordersCreator,
		stocksCreator: stocksCreator,
	}
}

func (u *TxManager[OrderRepoType, StocksRepoType]) WithinTransaction(ctx context.Context, f func(ctx context.Context, orders OrderRepoType, stocks StocksRepoType) error) error {
	tx, err := u.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	err = f(ctx, u.ordersCreator(tx), u.stocksCreator(tx))
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
