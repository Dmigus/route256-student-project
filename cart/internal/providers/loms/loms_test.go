//go:build unit

package loms

import (
	"context"
	"fmt"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"route256.ozon.ru/project/cart/internal/models"
	"testing"
)

func TestLOMSIsItemAvailablePositive(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx   context.Context
		skuId int64
		count uint16
	}
	tests := []struct {
		name      string
		args      args
		mockSetup func(testHelper)
		want      bool
	}{
		{
			name: "enough",
			args: args{
				ctx:   context.Background(),
				skuId: 123,
				count: 5,
			},
			mockSetup: func(h testHelper) {
				h.getNumberMock.Expect(minimock.AnyContext, 123).Return(6, nil)
			},
			want: true,
		},
		{
			name: "not enough",
			args: args{
				ctx:   context.Background(),
				skuId: 124,
				count: 5,
			},
			mockSetup: func(h testHelper) {
				h.getNumberMock.Expect(minimock.AnyContext, 124).Return(3, nil)
			},
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			helper := newTestHelper(t)
			tt.mockSetup(helper)
			ok, err := helper.loms.IsItemAvailable(tt.args.ctx, tt.args.skuId, tt.args.count)
			require.NoError(t, err)
			assert.Equal(t, tt.want, ok)
		})
	}
}

func TestLOMSIsItemAvailableError(t *testing.T) {
	t.Parallel()
	errorToThrow := fmt.Errorf("oops error")
	helper := newTestHelper(t)
	helper.getNumberMock.Expect(minimock.AnyContext, 123).Return(0, errorToThrow)
	_, err := helper.loms.IsItemAvailable(context.Background(), 123, 5)
	assert.ErrorIs(t, err, errorToThrow)
}

func TestLOMS_CreateOrder(t *testing.T) {
	t.Parallel()
	helper := newTestHelper(t)
	items := []models.CartItem{{124, 5}}
	helper.createOrderMock.Expect(minimock.AnyContext, 123, items).Return(1, nil)
	orderId, err := helper.loms.CreateOrder(context.Background(), 123, items)
	require.NoError(t, err)
	assert.Equal(t, int64(1), orderId)
}

func TestLOMS_CreateOrderEmpty(t *testing.T) {
	t.Parallel()
	helper := newTestHelper(t)
	items := []models.CartItem{}
	_, err := helper.loms.CreateOrder(context.Background(), 123, items)
	require.ErrorIs(t, err, errEmptyOrder)
}
