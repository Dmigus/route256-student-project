//go:build unit
// +build unit

package stocksinfogetter

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"testing"
)

type testHelper struct {
	stockRepoMock *mStockRepoMockGetNumOfAvailable
	getter        *StocksInfoGetter
}

func newTestHelper(t *testing.T) testHelper {
	mc := minimock.NewController(t)
	helper := testHelper{}
	stocks := NewStockRepoMock(mc)
	helper.stockRepoMock = &(stocks.GetNumOfAvailableMock)
	txM := NewTxManagerMock(mc)
	txM.WithinTransactionMock.Set(func(ctx context.Context, f1 func(ctx context.Context, anyVal any, stocks StockRepo) error) error {
		return f1(ctx, nil, stocks)
	})
	helper.getter = NewGetter(txM)
	return helper
}
