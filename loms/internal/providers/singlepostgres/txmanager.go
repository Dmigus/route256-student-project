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
	// TxManagerOne предоставляет обёртку для юзкейсов, позволяющую использовать один провайдер в атомарных сценариях
	TxManagerOne[T any] struct {
		conn    TxBeginner
		creator func(pgx.Tx) T
	}
	// TxManagerTwo предоставляет обёртку для юзкейсов, позволяющую использовать два провайдера в атомарных сценариях
	TxManagerTwo[T1, T2 any] struct {
		conn     TxBeginner
		creator1 func(pgx.Tx) T1
		creator2 func(pgx.Tx) T2
	}
)

// NewTxManagerOne создаёт новую обёртку. conn - объект, который будет использоваться для открытия новых транзакций. creator - функция создания экземпляра провайдера, привязанного к новой транзакции
func NewTxManagerOne[T any](conn TxBeginner, creator func(conn pgx.Tx) T) *TxManagerOne[T] {
	return &TxManagerOne[T]{
		conn:    conn,
		creator: creator,
	}
}

// NewTxManagerTwo создаёт новую обёртку. conn - объект, который будет использоваться для открытия новых транзакций. creator1 - функция создания экземпляра репозитория заказов, привязанного к новой транзакции. creator2 - функция создания экземпляра репозитория стоков, привязанного к новой транзакции.
func NewTxManagerTwo[T1, T2 any](conn TxBeginner, creator1 func(conn pgx.Tx) T1, creator2 func(conn pgx.Tx) T2) *TxManagerTwo[T1, T2] {
	return &TxManagerTwo[T1, T2]{
		conn:     conn,
		creator1: creator1,
		creator2: creator2,
	}
}

// WithinTransaction создаёт новую транзакию, связанные с ней репозитории заказов и стоков и выполняет функцию f с этими репозиториями. Если функция завершилась без ошибки, транзакция фиксируется, иначе откатывается.
func (u *TxManagerOne[T]) WithinTransaction(ctx context.Context, f func(ctx context.Context, provider T) error) error {
	tx, err := u.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	err = f(ctx, u.creator(tx))
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// WithinTransaction создаёт новую транзакию, связанные с ней репозитории заказов и стоков и выполняет функцию f с этими репозиториями. Если функция завершилась без ошибки, транзакция фиксируется, иначе откатывается.
func (u *TxManagerTwo[T1, T2]) WithinTransaction(ctx context.Context, f func(ctx context.Context, provider1 T1, provider2 T2) error) error {
	tx, err := u.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	err = f(ctx, u.creator1(tx), u.creator2(tx))
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
