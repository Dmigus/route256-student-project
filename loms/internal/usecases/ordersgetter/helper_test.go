//go:build unit
// +build unit

package ordersgetter

import (
	"github.com/gojuno/minimock/v3"
	"testing"
)

type testHelper struct {
	orderLoadRepoMock *mOrdersStorageMockLoad
	getter            *OrdersGetter
}

func newTestHelper(t *testing.T) testHelper {
	mc := minimock.NewController(t)
	helper := testHelper{}
	orders := NewOrdersStorageMock(mc)
	helper.orderLoadRepoMock = &(orders.LoadMock)
	helper.getter = NewOrdersGetter(orders)
	return helper
}
