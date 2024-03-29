// Package inmemory содержит in-memory реализацию хранилища стоков и заказов
package inmemory

import (
	"context"
)

// TxManager предоставляет обёртку для юзкейсов, позволяющую использовать in-memory хранилища в атомарных сценариях
type TxManager[OrderRepoType, StocksRepoType any] struct {
	orders OrderRepoType
	stocks StocksRepoType
}

// NewTxManager создаёт новый экземпляр TxManager
func NewTxManager[OrderRepoType, StocksRepoType any](orders OrderRepoType, stocks StocksRepoType) *TxManager[OrderRepoType, StocksRepoType] {
	return &TxManager[OrderRepoType, StocksRepoType]{
		orders: orders,
		stocks: stocks,
	}
}

// WithinTransaction исполняет функцию f с контекстом ctx и объектами репозяитория, переданными ему в конструкторе
func (u *TxManager[OrderRepoType, StocksRepoType]) WithinTransaction(ctx context.Context, f func(ctx context.Context, orders OrderRepoType, stocks StocksRepoType) error) error {
	return f(ctx, u.orders, u.stocks)
}
