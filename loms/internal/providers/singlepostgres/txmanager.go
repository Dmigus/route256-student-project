package singlepostgres

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type (
	TxBeginner2 interface {
		Begin(ctx context.Context) (pgx.Tx, error)
	}
	TxManager[OrderRepoType, StocksRepoType any] struct {
		conn          TxBeginner2
		ordersCreator func(conn TxBeginner2) OrderRepoType
		stocksCreator func(conn TxBeginner2) StocksRepoType
	}
)

func NewTxManager[OrderRepoType, StocksRepoType any](conn TxBeginner2, ordersCreator func(conn TxBeginner2) OrderRepoType, stocksCreator func(conn TxBeginner2) StocksRepoType) *TxManager[OrderRepoType, StocksRepoType] {
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
