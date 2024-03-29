//go:build unit
// +build unit

package ordersgetter

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"testing"
)

type testHelper struct {
	orderLoadRepoMock *mOrderRepoMockLoad
	getter            *OrdersGetter
}

func newTestHelper(t *testing.T) testHelper {
	mc := minimock.NewController(t)
	helper := testHelper{}
	ordersMock := NewOrderRepoMock(mc)
	helper.orderLoadRepoMock = &(ordersMock.LoadMock)
	txM := NewTxManagerMock(mc)
	txM.WithinTransactionMock.Set(func(ctx context.Context, f1 func(ctx context.Context, orders OrderRepo, _ any) error) error {
		return f1(ctx, ordersMock, nil)
	})
	helper.getter = NewOrdersGetter(txM)
	return helper
}
