//go:build unit
// +build unit

package stocksinfogetter

import (
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
	helper.getter = NewGetter(stocks)
	return helper
}
