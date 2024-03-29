package stocksinfogetter

import (
	"context"
	"fmt"
)

type (
	StockRepo interface {
		GetNumOfAvailable(context.Context, int64) (uint64, error)
	}
	txManager interface {
		WithinTransaction(context.Context, func(ctx context.Context, _ any, stocks StockRepo) error) error
	}
	StocksInfoGetter struct {
		tx txManager
	}
)

func NewGetter(tx txManager) *StocksInfoGetter {
	return &StocksInfoGetter{tx: tx}
}

func (g *StocksInfoGetter) GetNumOfAvailable(ctx context.Context, skuId int64) (cnt uint64, err error) {
	err = g.tx.WithinTransaction(ctx, func(ctx context.Context, _ any, stocks StockRepo) error {
		cnt, err = stocks.GetNumOfAvailable(ctx, skuId)
		if err != nil {
			return fmt.Errorf("could not get number of available points for item %d: %w", skuId, err)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return cnt, nil
}
