//go:build unit

package checkouter

import (
	"context"
	"fmt"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"route256.ozon.ru/project/cart/internal/models"
	"testing"
)

func TestCheckouterPositive(t *testing.T) {
	t.Parallel()
	helper := newTestHelper(t)
	cart := models.NewCart()
	ctx := context.Background()
	cart.Add(ctx, 456, 10)
	cart.Add(ctx, 789, 5)
	helper.getCartRepoMock.Expect(minimock.AnyContext, 123).Return(cart, nil)
	cartItems := cart.ListItemsSorted(ctx)
	helper.createOrderMock.Expect(minimock.AnyContext, 123, cartItems).Return(5, nil)
	helper.clearCartRepoMock.Expect(minimock.AnyContext, 123).Return()
	orderId, err := helper.service.Checkout(ctx, 123)
	require.NoError(t, err)
	assert.Equal(t, int64(5), orderId)
}

func TestCheckouterErrors(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx    context.Context
		userId int64
	}
	errorToThrow := fmt.Errorf("oops error")
	tests := []struct {
		name      string
		mockSetup func(testHelper)
		args      args
	}{
		{
			name: "error getting cart",
			mockSetup: func(h testHelper) {
				h.getCartRepoMock.Expect(minimock.AnyContext, 123).Return(nil, errorToThrow)
			},
			args: args{
				ctx:    context.Background(),
				userId: 123,
			},
		},
		{
			name: "error creating order",
			mockSetup: func(h testHelper) {
				cart := models.NewCart()
				h.getCartRepoMock.Expect(minimock.AnyContext, 123).Return(cart, nil)
				items := cart.ListItems(context.Background())
				h.createOrderMock.Expect(minimock.AnyContext, 123, items).Return(0, errorToThrow)
			},
			args: args{
				ctx:    context.Background(),
				userId: 123,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			helper := newTestHelper(t)
			tt.mockSetup(helper)
			_, err := helper.service.Checkout(tt.args.ctx, tt.args.userId)
			assert.ErrorIs(t, err, errorToThrow)
		})
	}
}
