//go:build unit
// +build unit

package orderscreator

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"testing"
)

type testHelper struct {
	stocksMock          *mStockRepoMockReserve
	orderCreateRepoMock *mOrderRepoMockCreate
	orderSaveRepoMock   *mOrderRepoMockSave
	creator             *OrdersCreator
}

func newTestHelper(t *testing.T) testHelper {
	mc := minimock.NewController(t)
	helper := testHelper{}
	orders := NewOrderRepoMock(mc)
	stocks := NewStockRepoMock(mc)
	helper.orderSaveRepoMock = &(orders.SaveMock)
	helper.stocksMock = &(stocks.ReserveMock)
	helper.orderCreateRepoMock = &(orders.CreateMock)
	txM := NewTxManagerMock(mc)
	txM.WithinTransactionMock.Set(func(ctx context.Context, f1 func(ctx context.Context, orders OrderRepo, stocks StockRepo) error) (err error) {
		return f1(ctx, orders, stocks)
	})
	helper.creator = NewOrdersCreator(txM)
	return helper
}
