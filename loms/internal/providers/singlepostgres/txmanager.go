// Package singlepostgres содержит реализацию хранилища стоков и заказов в одной БД PostgreSQL
package singlepostgres

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type (
	// TxBeginner объект, умеющий открывать новую транзакцию
	TxBeginner interface {
		Begin(ctx context.Context) (pgx.Tx, error)
	}
	// TxManager предоставляет обёртку для юзкейсов, позволяющую использовать репозитории в атомарных сценариях
	TxManager[OrderRepoType, StocksRepoType any] struct {
		conn          TxBeginner
		ordersCreator func(pgx.Tx) OrderRepoType
		stocksCreator func(pgx.Tx) StocksRepoType
	}
)

// NewTxManager создаёт новую обёртку. conn - объект, который будет использоваться для открытия новых транзакций. ordersCreator - функция создания экземпляра репозитория заказов, привязанного к новой транзакции. stocksCreator - функция создания экземпляра репозитория стоков, привязанного к новой транзакции.
func NewTxManager[OrderRepoType, StocksRepoType any](conn TxBeginner, ordersCreator func(conn pgx.Tx) OrderRepoType, stocksCreator func(conn pgx.Tx) StocksRepoType) *TxManager[OrderRepoType, StocksRepoType] {
	return &TxManager[OrderRepoType, StocksRepoType]{
		conn:          conn,
		ordersCreator: ordersCreator,
		stocksCreator: stocksCreator,
	}
}

// WithinTransaction создаёт новую транзакию, связанные с ней репозитории заказов и стоков и выполняет функцию f с этими репозиториями. Если функция завершилась без ошибки, транзакция фиксируется, иначе откатывается.
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
