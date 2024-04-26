package multipostgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"route256.ozon.ru/project/loms/internal/providers/multipostgres/shardmanager"
)

type (
	TransactionCreator interface {
		GetTransaction(ctx context.Context, beginner TxBeginner) (pgx.Tx, error)
	}
	TxManager1[T any] struct {
		creator func(TransactionCreator) T
	}
	TxManager2[T1 any, T2 any] struct {
		creator1 func(TransactionCreator) T1
		creator2 func(TransactionCreator) T2
	}
	TxManager3[T1 any, T2 any, T3 any] struct {
		creator1 func(TransactionCreator) T1
		creator2 func(TransactionCreator) T2
		creator3 func(TransactionCreator) T3
	}
)

func NewTxManager1[T any](createProviderFunc func(TransactionCreator) T) *TxManager1[T] {
	return &TxManager1[T]{
		creator: createProviderFunc,
	}
}

func NewTxManager2[T1 any, T2 any](creator1 func(TransactionCreator) T1, creator2 func(TransactionCreator) T2) *TxManager2[T1, T2] {
	return &TxManager2[T1, T2]{
		creator1: creator1,
		creator2: creator2,
	}
}

func NewTxManager3[T1 any, T2 any, T3 any](creator1 func(TransactionCreator) T1, creator2 func(TransactionCreator) T2, creator3 func(creator TransactionCreator) T3) *TxManager3[T1, T2, T3] {
	return &TxManager3[T1, T2, T3]{
		creator1: creator1,
		creator2: creator2,
		creator3: creator3,
	}
}

func (txm *TxManager1[T]) WithinTransaction(ctx context.Context, f func(ctx context.Context, provider T) bool) error {
	tr := newDistributedTransaction()
	command := &command1[T]{
		provider: txm.creator(tr),
		f:        f,
	}
	return tr.WithinTransaction(ctx, command)
}

func (txm *TxManager2[T1, T2]) WithinTransaction(ctx context.Context, f func(ctx context.Context, provider1 T1, provider2 T2) bool) error {
	tr := newDistributedTransaction()
	command := &command2[T1, T2]{
		provider1: txm.creator1(tr),
		provider2: txm.creator2(tr),
		f:         f,
	}
	return tr.WithinTransaction(ctx, command)
}

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

func OrderIDToShardKey(orderID int64) shardmanager.ShardKey {
	return shardmanager.ShardKey(orderID % 1000)
}
