package singlepostres

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

func InTx(conn txBeginner, ctx context.Context, txOptions pgx.TxOptions, f wrapableFunc) error {
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
