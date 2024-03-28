//go:build unit
// +build unit

package orderscreator

import (
	"github.com/gojuno/minimock/v3"
	"testing"
)

type testHelper struct {
	stocksMock          *mStocksStorageMockReserve
	orderCreateRepoMock *mOrdersStorageMockCreate
	orderSaveRepoMock   *mOrdersStorageMockSave
	creator             *OrdersCreator
}

func newTestHelper(t *testing.T) testHelper {
	mc := minimock.NewController(t)
	helper := testHelper{}
	orders := NewOrdersStorageMock(mc)
	stocks := NewStocksStorageMock(mc)
	helper.orderSaveRepoMock = &(orders.SaveMock)
	helper.stocksMock = &(stocks.ReserveMock)
	helper.orderCreateRepoMock = &(orders.CreateMock)
	helper.creator = NewOrdersCreator(orders, stocks)
	return helper
}
