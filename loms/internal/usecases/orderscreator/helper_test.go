//go:build unit
// +build unit

package orderscreator

import (
	"github.com/gojuno/minimock/v3"
	"testing"
)

type testHelper struct {
	stocksMock        *mStocksStorageMockReserve
	generatorMock     *mOrderIdGeneratorMockNewId
	orderSaveRepoMock *mOrdersStorageMockSave
	creator           *OrdersCreator
}

func newTestHelper(t *testing.T) testHelper {
	mc := minimock.NewController(t)
	helper := testHelper{}
	orders := NewOrdersStorageMock(mc)
	generator := NewOrderIdGeneratorMock(mc)
	stocks := NewStocksStorageMock(mc)
	helper.orderSaveRepoMock = &(orders.SaveMock)
	helper.stocksMock = &(stocks.ReserveMock)
	helper.generatorMock = &(generator.NewIdMock)
	helper.creator = NewOrdersCreator(generator, orders, stocks)
	return helper
}
