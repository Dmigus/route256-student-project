//go:build unit
// +build unit

package orderspayer

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"testing"
)

type testHelper struct {
	orderLoadRepoMock *mOrderRepoMockLoad
	orderSaveRepoMock *mOrderRepoMockSave
	stocksRepoMock    *mStockRepoMockRemoveReserved
	payer             *OrdersPayer
}

func newTestHelper(t *testing.T) testHelper {
	mc := minimock.NewController(t)
	helper := testHelper{}
	orders := NewOrderRepoMock(mc)
	stocks := NewStockRepoMock(mc)
	helper.orderLoadRepoMock = &(orders.LoadMock)
	helper.orderSaveRepoMock = &(orders.SaveMock)
	helper.stocksRepoMock = &(stocks.RemoveReservedMock)
	txM := NewTxManagerMock(t)
	txM.WithinTransactionMock.Set(func(ctx context.Context, f1 func(ctx context.Context, orders OrderRepo, stocks StockRepo) error) error {
		return f1(ctx, orders, stocks)
	})
	helper.payer = NewOrdersPayer(txM)
	return helper
}
