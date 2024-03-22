//go:build unit

package orderscanceller

import (
	"context"
	"fmt"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"route256.ozon.ru/project/loms/internal/models"
	"testing"
)

func TestOrderCancellerPositive(t *testing.T) {
	t.Parallel()
	helper := newTestHelper(t)
	order := models.NewOrder(123, 1234)
	order.Status = models.AwaitingPayment
	order.IsItemsReserved = true
	items := []models.OrderItem{{34, 10}}
	order.Items = items
	helper.orderLoadRepoMock.Expect(minimock.AnyContext, 1234).Return(order, nil)
	helper.stocksMock.Expect(minimock.AnyContext, items).Return(nil)
	helper.orderSaveRepoMock.Expect(minimock.AnyContext, order).Return(nil)
	err := helper.canceller.Cancel(minimock.AnyContext, 1234)
	require.NoError(t, err)
	assert.Equal(t, models.Cancelled, order.Status)
	assert.False(t, order.IsItemsReserved)
}

func TestOrderCancellerErrors(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx     context.Context
		orderId int64
	}
	errorToThrow := fmt.Errorf("oops error")
	tests := []struct {
		name      string
		mockSetup func(testHelper)
		args      args
		err       error
	}{
		{
			name: "error getting order",
			mockSetup: func(helper testHelper) {
				helper.orderLoadRepoMock.Expect(minimock.AnyContext, 1234).Return(nil, errorToThrow)
			},
			args: args{
				ctx:     context.Background(),
				orderId: 1234,
			},
			err: errorToThrow,
		},
		{
			name: "error wrong status",
			mockSetup: func(helper testHelper) {
				order := models.NewOrder(123, 1234)
				order.Status = models.Cancelled
				helper.orderLoadRepoMock.Expect(minimock.AnyContext, 1234).Return(order, nil)
			},
			args: args{
				ctx:     context.Background(),
				orderId: 1234,
			},
			err: models.ErrWrongOrderStatus,
		},
		{
			name: "error canceling reserved stocks",
			mockSetup: func(helper testHelper) {
				order := models.NewOrder(123, 1234)
				order.Status = models.AwaitingPayment
				order.IsItemsReserved = true
				items := make([]models.OrderItem, 0)
				order.Items = items
				helper.orderLoadRepoMock.Expect(minimock.AnyContext, 1234).Return(order, nil)
				helper.stocksMock.Expect(minimock.AnyContext, items).Return(errorToThrow)
			},
			args: args{
				ctx:     context.Background(),
				orderId: 1234,
			},
			err: errorToThrow,
		},
		{
			name: "error adding payed stocks",
			mockSetup: func(helper testHelper) {
				order := models.NewOrder(123, 1234)
				order.Status = models.Payed
				items := []models.OrderItem{{12, 6}}
				order.Items = items
				helper.orderLoadRepoMock.Expect(minimock.AnyContext, 1234).Return(order, nil)
				helper.addItemsMock.Expect(minimock.AnyContext, items).Return(errorToThrow)
			},
			args: args{
				ctx:     context.Background(),
				orderId: 1234,
			},
			err: errorToThrow,
		},
		{
			name: "error saving order",
			mockSetup: func(helper testHelper) {
				order := models.NewOrder(123, 1234)
				order.Status = models.AwaitingPayment
				order.IsItemsReserved = true
				items := make([]models.OrderItem, 0)
				order.Items = items
				helper.orderLoadRepoMock.Expect(minimock.AnyContext, 1234).Return(order, nil)
				helper.stocksMock.Expect(minimock.AnyContext, items).Return(nil)
				helper.orderSaveRepoMock.Expect(minimock.AnyContext, order).Return(errorToThrow)
			},
			args: args{
				ctx:     context.Background(),
				orderId: 1234,
			},
			err: errorToThrow,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			helper := newTestHelper(t)
			tt.mockSetup(helper)
			err := helper.canceller.Cancel(tt.args.ctx, tt.args.orderId)
			assert.ErrorIs(t, err, tt.err)
		})
	}
}
