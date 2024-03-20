//go:build unit

package loms

import (
	"github.com/gojuno/minimock/v3"
	"testing"
)

type testHelper struct {
	createOrderMock *mLomsClientMockCreateOrder
	getNumberMock   *mLomsClientMockGetNumberOfItemInStocks
	loms            *LOMS
}

func newTestHelper(t *testing.T) testHelper {
	mc := minimock.NewController(t)
	helper := testHelper{}
	client := NewLomsClientMock(mc)
	helper.getNumberMock = &(client.GetNumberOfItemInStocksMock)
	helper.createOrderMock = &(client.CreateOrderMock)
	helper.loms = NewLOMSProvider(client)
	return helper
}
