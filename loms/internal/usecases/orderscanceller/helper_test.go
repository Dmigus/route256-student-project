//go:build unit
// +build unit

package orderscanceller

import (
	"github.com/gojuno/minimock/v3"
	"testing"
)

type testHelper struct {
	stocksMock        *mStockCancellerMockCancelReserved
	orderLoadRepoMock *mOrderRepoMockLoad
	orderSaveRepoMock *mOrderRepoMockSave
	canceller         *OrderCanceller
}

func newTestHelper(t *testing.T) testHelper {
	mc := minimock.NewController(t)
	helper := testHelper{}
	orders := NewOrderRepoMock(mc)
	stocks := NewStockCancellerMock(mc)
	helper.orderLoadRepoMock = &(orders.LoadMock)
	helper.orderSaveRepoMock = &(orders.SaveMock)
	helper.stocksMock = &(stocks.CancelReservedMock)
	helper.canceller = NewOrderCanceller(orders, stocks)
	return helper
}
