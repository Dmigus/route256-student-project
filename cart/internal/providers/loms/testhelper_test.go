//go:build unit

package loms

import (
	"github.com/gojuno/minimock/v3"
	"testing"
)

type testHelper struct {
	createOrderMock *mLomsClientMockOrderCreate
	getNumberMock   *mLomsClientMockGetNumberOfItemInStocks
	loms            *LOMS
}

func newTestHelper(t *testing.T) testHelper {
	mc := minimock.NewController(t)
	helper := testHelper{}
	client := NewLomsClientMock(mc)
	helper.getNumberMock = &(client.GetNumberOfItemInStocksMock)
	helper.createOrderMock = &(client.OrderCreateMock)
	helper.loms = NewLOMSProvider(client)
	return helper
}
