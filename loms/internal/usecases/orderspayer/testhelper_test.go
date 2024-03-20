//go:build unit
// +build unit

package orderspayer

import (
	"github.com/gojuno/minimock/v3"
	"testing"
)

type testHelper struct {
	orderLoadRepoMock *mOrderRepoMockLoad
	orderSaveRepoMock *mOrderRepoMockSave
	stocksRepoMock    *mStockRemoverMockRemoveReserved
	payer             *OrdersPayer
}

func newTestHelper(t *testing.T) testHelper {
	mc := minimock.NewController(t)
	helper := testHelper{}
	orders := NewOrderRepoMock(mc)
	stocks := NewStockRemoverMock(mc)
	helper.orderLoadRepoMock = &(orders.LoadMock)
	helper.orderSaveRepoMock = &(orders.SaveMock)
	helper.stocksRepoMock = &(stocks.RemoveReservedMock)
	helper.payer = NewOrdersPayer(orders, stocks)
	return helper
}
