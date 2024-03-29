//go:build unit
// +build unit

package orderscanceller

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"testing"
)

type testHelper struct {
	stocksMock        *mStockRepoMockCancelReserved
	addItemsMock      *mStockRepoMockAddItems
	orderLoadRepoMock *mOrderRepoMockLoad
	orderSaveRepoMock *mOrderRepoMockSave
	canceller         *OrderCanceller
}

func newTestHelper(t *testing.T) testHelper {
	mc := minimock.NewController(t)
	helper := testHelper{}
	orders := NewOrderRepoMock(mc)
	stocks := NewStockRepoMock(mc)
	helper.orderLoadRepoMock = &(orders.LoadMock)
	helper.orderSaveRepoMock = &(orders.SaveMock)
	helper.stocksMock = &(stocks.CancelReservedMock)
	helper.addItemsMock = &(stocks.AddItemsMock)
	trM := NewTrManagerMock(t)
	trM.WithinTransactionMock.Set(func(ctx context.Context, f1 func(ctx context.Context, orders OrderRepo, stocks StockRepo) error) (err error) {
		return f1(ctx, orders, stocks)
	})
	helper.canceller = NewOrderCanceller(trM)
	return helper
}
