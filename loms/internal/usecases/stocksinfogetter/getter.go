package stocksinfogetter

import "context"

type stockRepo interface {
	GetNumOfAvailable(context.Context, uint32) (uint64, error)
}

type Getter struct {
	stocks stockRepo
}

func (g *Getter) GetNumOfAvailable(ctx context.Context, skuId uint32) (uint64, error) {
	return g.stocks.GetNumOfAvailable(ctx, skuId)
}
