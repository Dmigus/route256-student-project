// Package singlepostres предназначен для реализации репозиториев стоков и заказов на основе одной БД в PostgresSQL.
package singlepostgres

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type (
	exclusiveType int
	wrapableFunc  func(ctx context.Context) error
	txBeginner    interface {
		BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	}
)

const trKey = exclusiveType(0)

// InTx выполняет функцию f в транзакции, созданной из conn с парамерами txOptions. Если f вернула err != nil, то происходит Rollback, иначе Commit
func InTx(ctx context.Context, conn txBeginner, txOptions pgx.TxOptions, f wrapableFunc) error {
	tx, err := conn.BeginTx(ctx, txOptions)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	txCtx := context.WithValue(ctx, trKey, tx)
	err = f(txCtx)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
