package inmemory

import (
	"context"
)

type TxManager[OrderRepoType, StocksRepoType any] struct {
	orders OrderRepoType
	stocks StocksRepoType
}

func NewTxManager[OrderRepoType, StocksRepoType any](orders OrderRepoType, stocks StocksRepoType) *TxManager[OrderRepoType, StocksRepoType] {
	return &TxManager[OrderRepoType, StocksRepoType]{
		orders: orders,
		stocks: stocks,
	}
}

func (u *TxManager[OrderRepoType, StocksRepoType]) WithinTransaction(ctx context.Context, f func(ctx context.Context, orders OrderRepoType, stocks StocksRepoType) error) error {
	return f(ctx, u.orders, u.stocks)
}
