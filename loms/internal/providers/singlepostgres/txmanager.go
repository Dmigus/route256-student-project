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
	// TxManagerThree предоставляет обёртку для юзкейсов, позволяющую использовать три провайдера в атомарных сценариях
	TxManagerThree[T1, T2, T3 any] struct {
		conn     TxBeginner
		creator1 func(pgx.Tx) T1
		creator2 func(pgx.Tx) T2
		creator3 func(pgx.Tx) T3
	}
)

// NewTxManagerOne создаёт новую обёртку. conn - объект, который будет использоваться для открытия новых транзакций. creator - функция создания экземпляра провайдера, привязанного к новой транзакции
func NewTxManagerOne[T any](conn TxBeginner, creator func(conn pgx.Tx) T) *TxManagerOne[T] {
	return &TxManagerOne[T]{
		conn:    conn,
		creator: creator,
	}
}

// NewTxManagerTwo создаёт новую обёртку. conn - объект, который будет использоваться для открытия новых транзакций. creator'ы - функции создания экземпляра провайдера, привязанного к новой транзакции.
func NewTxManagerTwo[T1, T2 any](conn TxBeginner, creator1 func(conn pgx.Tx) T1, creator2 func(conn pgx.Tx) T2) *TxManagerTwo[T1, T2] {
	return &TxManagerTwo[T1, T2]{
		conn:     conn,
		creator1: creator1,
		creator2: creator2,
	}
}

// NewTxManagerThree создаёт новую обёртку. conn - объект, который будет использоваться для открытия новых транзакций. creator'ы - функции создания экземпляра провайдера, привязанного к новой транзакции.
func NewTxManagerThree[T1, T2, T3 any](conn TxBeginner, creator1 func(conn pgx.Tx) T1, creator2 func(conn pgx.Tx) T2, creator3 func(conn pgx.Tx) T3) *TxManagerThree[T1, T2, T3] {
	return &TxManagerThree[T1, T2, T3]{
		conn:     conn,
		creator1: creator1,
		creator2: creator2,
		creator3: creator3,
	}
}

// WithinTransaction создаёт новую транзакию, связанные с ней репозитории заказов и стоков и выполняет функцию f с этими провайдерами. Если функция завершилась без ошибки, транзакция фиксируется, иначе откатывается.
func (u *TxManagerOne[T]) WithinTransaction(ctx context.Context, f func(ctx context.Context, provider T) bool) error {
	tx, err := u.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if f(ctx, u.creator(tx)) {
		return tx.Commit(ctx)
	}
	return nil
}

// WithinTransaction создаёт новую транзакию, связанные с ней репозитории заказов и стоков и выполняет функцию f с двумя провайдерами. Если функция завершилась без ошибки, транзакция фиксируется, иначе откатывается.
func (u *TxManagerTwo[T1, T2]) WithinTransaction(ctx context.Context, f func(ctx context.Context, provider1 T1, provider2 T2) bool) error {
	tx, err := u.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if f(ctx, u.creator1(tx), u.creator2(tx)) {
		return tx.Commit(ctx)
	}
	return nil
}

// WithinTransaction создаёт новую транзакию, связанные с ней репозитории заказов и стоков и выполняет функцию f с тремя провайдерами. Если функция завершилась без ошибки, транзакция фиксируется, иначе откатывается.
func (u *TxManagerThree[T1, T2, T3]) WithinTransaction(ctx context.Context, f func(ctx context.Context, provider1 T1, provider2 T2, provider3 T3) bool) error {
	tx, err := u.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if f(ctx, u.creator1(tx), u.creator2(tx), u.creator3(tx)) {
		return tx.Commit(ctx)
	}
	return nil
}
