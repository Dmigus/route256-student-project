package stocksinfogetter

import "context"

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
	return g.stocks.GetNumOfAvailable(ctx, skuId)
}
