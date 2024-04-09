// Package inmemory содержит in-memory реализацию хранилища стоков и заказов
package inmemory

import (
	"context"
)

// TxManagerOne предоставляет обёртку для юзкейсов, позволяющую использовать один in-memory провайдер в атомарных сценариях
type TxManagerOne[T any] struct {
	provider T
}

// TxManagerTwo предоставляет обёртку для юзкейсов, позволяющую использовать два in-memory провайдера в атомарных сценариях
type TxManagerTwo[T1, T2 any] struct {
	provider1 T1
	provider2 T2
}

// NewTxManagerOne создаёт новый экземпляр TxManagerOne
func NewTxManagerOne[T any](provider T) *TxManagerOne[T] {
	return &TxManagerOne[T]{
		provider: provider,
	}
}

// NewTxManagerTwo создаёт новый экземпляр TxManagerTwo
func NewTxManagerTwo[T1, T2 any](provider1 T1, provider2 T2) *TxManagerTwo[T1, T2] {
	return &TxManagerTwo[T1, T2]{
		provider1: provider1,
		provider2: provider2,
	}
}

// WithinTransaction исполняет функцию f с контекстом ctx и одним провайдером, переданными ему в конструкторе
func (u *TxManagerOne[T]) WithinTransaction(ctx context.Context, f func(ctx context.Context, provider T) error) error {
	return f(ctx, u.provider)
}

// WithinTransaction исполняет функцию f с контекстом ctx и двумя провайдерами, переданными ему в конструкторе
func (u *TxManagerTwo[T1, T2]) WithinTransaction(ctx context.Context, f func(ctx context.Context, provider1 T1, provider2 T2) error) error {
	return f(ctx, u.provider1, u.provider2)
}
