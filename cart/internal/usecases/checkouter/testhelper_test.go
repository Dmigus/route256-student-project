//go:build unit
// +build unit

package checkouter

import (
	"github.com/gojuno/minimock/v3"
	"testing"
)

type testHelper struct {
	getCartRepoMock   *mRepositoryMockGetCart
	clearCartRepoMock *mRepositoryMockClearCartReliable
	createOrderMock   *mOrderSystemMockCreateOrder
	service           *Checkouter
}

func newTestHelper(t *testing.T) testHelper {
	mc := minimock.NewController(t)
	helper := testHelper{}
	repo := NewRepositoryMock(mc)
	orders := NewOrderSystemMock(mc)
	helper.getCartRepoMock = &(repo.GetCartMock)
	helper.clearCartRepoMock = &(repo.ClearCartReliableMock)
	helper.createOrderMock = &(orders.CreateOrderMock)
	helper.service = NewCheckouter(repo, orders)
	return helper
}
