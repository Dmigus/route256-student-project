package stocksinfogetter

import (
	"context"
	"fmt"
)

type stockRepo interface {
	GetNumOfAvailable(context.Context, int64) (uint64, error)
}

type StocksInfoGetter struct {
	stocks stockRepo
}

func NewGetter(stocks stockRepo) *StocksInfoGetter {
	return &StocksInfoGetter{stocks: stocks}
}

func (g *StocksInfoGetter) GetNumOfAvailable(ctx context.Context, skuId int64) (uint64, error) {
	num, err := g.stocks.GetNumOfAvailable(ctx, skuId)
	if err != nil {
		return 0, fmt.Errorf("could not get number of available points for item %d: %w", skuId, err)
	}
	return num, nil
}
