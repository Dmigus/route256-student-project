// Package multipostgres содержит реализацию работы с транзакцией, распределённой по нескольким узлам Postgres
package multipostgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"route256.ozon.ru/project/loms/internal/providers/multipostgres/shardmanager"
)

type (
	// TransactionCreator это объект, который умеет выдавать транзакцию из beginner. Реализуется типом distributedTransaction
	TransactionCreator interface {
		GetTransaction(ctx context.Context, beginner TxBeginner) (pgx.Tx, error)
	}
	// TxManager1 это менеждер, предназначенный для работы с одним провайдером
	TxManager1[T any] struct {
		creator func(TransactionCreator) T
	}
	// TxManager2 это менеждер, предназначенный для работы с двумя провайдерами
	TxManager2[T1 any, T2 any] struct {
		creator1 func(TransactionCreator) T1
		creator2 func(TransactionCreator) T2
	}
	// TxManager3 это менеждер, предназначенный для работы с тремя провайдерами
	TxManager3[T1 any, T2 any, T3 any] struct {
		creator1 func(TransactionCreator) T1
		creator2 func(TransactionCreator) T2
		creator3 func(TransactionCreator) T3
	}
)

// NewTxManager1 создаёт новый TxManager1
func NewTxManager1[T any](createProviderFunc func(TransactionCreator) T) *TxManager1[T] {
	return &TxManager1[T]{
		creator: createProviderFunc,
	}
}

// NewTxManager2 создаёт новый TxManager2
func NewTxManager2[T1 any, T2 any](creator1 func(TransactionCreator) T1, creator2 func(TransactionCreator) T2) *TxManager2[T1, T2] {
	return &TxManager2[T1, T2]{
		creator1: creator1,
		creator2: creator2,
	}
}

// NewTxManager3 создаёт новый TxManager3
func NewTxManager3[T1 any, T2 any, T3 any](creator1 func(TransactionCreator) T1, creator2 func(TransactionCreator) T2, creator3 func(creator TransactionCreator) T3) *TxManager3[T1, T2, T3] {
	return &TxManager3[T1, T2, T3]{
		creator1: creator1,
		creator2: creator2,
		creator3: creator3,
	}
}

// WithinTransaction исполняет функцию f в одной распределённой транзакции
func (txm *TxManager1[T]) WithinTransaction(ctx context.Context, f func(ctx context.Context, provider T) bool) error {
	tr := newDistributedTransaction()
	command := &command1[T]{
		provider: txm.creator(tr),
		f:        f,
	}
	return tr.WithinTransaction(ctx, command)
}

// WithinTransaction исполняет функцию f в одной распределённой транзакции
func (txm *TxManager2[T1, T2]) WithinTransaction(ctx context.Context, f func(ctx context.Context, provider1 T1, provider2 T2) bool) error {
	tr := newDistributedTransaction()
	command := &command2[T1, T2]{
		provider1: txm.creator1(tr),
		provider2: txm.creator2(tr),
		f:         f,
	}
	return tr.WithinTransaction(ctx, command)
}

// WithinTransaction исполняет функцию f в одной распределённой транзакции
func (txm *TxManager3[T1, T2, T3]) WithinTransaction(ctx context.Context, f func(ctx context.Context, provider1 T1, provider2 T2, provider3 T3) bool) error {
	tr := newDistributedTransaction()
	command := &command3[T1, T2, T3]{
		provider1: txm.creator1(tr),
		provider2: txm.creator2(tr),
		provider3: txm.creator3(tr),
		f:         f,
	}
	return tr.WithinTransaction(ctx, command)
}

// OrderIDToShardBucket вычисляет номер бакета, к которому принадлежит заказ с id = orderID
func OrderIDToShardBucket(orderID int64) shardmanager.ShardBucket {
	return shardmanager.ShardBucket(orderID % shardmanager.BucketsNum)
}
