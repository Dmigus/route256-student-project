//go:build unit
// +build unit

package adder

import (
	"context"
	"fmt"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"route256.ozon.ru/project/cart/internal/models"
	"testing"
)

func TestAddItemWithoutErr(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx   context.Context
		user  int64
		skuId int64
		count uint16
	}
	tests := []struct {
		name      string
		args      args
		mockSetup func(testHelper)
	}{
		{
			name: "positive",
			args: args{
				ctx:   context.Background(),
				user:  123,
				skuId: 123,
				count: 1,
			},
			mockSetup: func(h testHelper) {
				cart := models.NewCart()
				h.productServiceMock.Expect(minimock.AnyContext, 123).Return(true, nil)
				h.stocksCheckerMock.Expect(minimock.AnyContext, 123, 1).Return(true, nil)
				h.getCartRepoMock.Expect(minimock.AnyContext, 123).Return(cart, nil)
				h.saveCartRepoMock.Expect(minimock.AnyContext, 123, cart).Return(nil)
			},
		},
		{
			name: "item not exists in product service",
			args: args{
				ctx:   context.Background(),
				user:  124,
				skuId: 124,
				count: 1,
			},
			mockSetup: func(h testHelper) {
				h.productServiceMock.Expect(minimock.AnyContext, 124).Return(false, nil)
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			helper := newTestHelper(t)
			tt.mockSetup(helper)
			err := helper.service.AddItem(tt.args.ctx, tt.args.user, tt.args.skuId, tt.args.count)
			assert.NoError(t, err, "must be no error")
		})
	}
}

func TestAddItemWithErrs(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx   context.Context
		user  int64
		skuId int64
		count uint16
	}
	errorToThrow := fmt.Errorf("oops error")
	tests := []struct {
		name      string
		mockSetup func(testHelper)
		args      args
		err       error
	}{
		{
			name: "error checking good presence",
			mockSetup: func(h testHelper) {
				h.productServiceMock.Expect(minimock.AnyContext, 123).Return(false, errorToThrow)
			},
			args: args{
				ctx:   context.Background(),
				user:  123,
				skuId: 123,
				count: 1,
			},
			err: errorToThrow,
		},
		{
			name: "error checking number of items",
			mockSetup: func(h testHelper) {
				h.productServiceMock.Expect(minimock.AnyContext, 123).Return(true, nil)
				h.stocksCheckerMock.Expect(minimock.AnyContext, 123, 1).Return(false, errorToThrow)
			},
			args: args{
				ctx:   context.Background(),
				user:  123,
				skuId: 123,
				count: 1,
			},
			err: errorToThrow,
		},
		{
			name: "error not enough items in stocks",
			mockSetup: func(h testHelper) {
				h.productServiceMock.Expect(minimock.AnyContext, 123).Return(true, nil)
				h.stocksCheckerMock.Expect(minimock.AnyContext, 123, 1).Return(false, nil)
			},
			args: args{
				ctx:   context.Background(),
				user:  123,
				skuId: 123,
				count: 1,
			},
			err: errNotEnoughNumInStocks,
		},
		{
			name: "error getting user cart",
			mockSetup: func(h testHelper) {
				h.productServiceMock.Expect(minimock.AnyContext, 123).Return(true, nil)
				h.stocksCheckerMock.Expect(minimock.AnyContext, 123, 1).Return(true, nil)
				h.getCartRepoMock.Expect(minimock.AnyContext, 123).Return(nil, errorToThrow)
			},
			args: args{
				ctx:   context.Background(),
				user:  123,
				skuId: 123,
				count: 1,
			},
			err: errorToThrow,
		},
		{
			name: "error saving user cart",
			mockSetup: func(h testHelper) {
				cart := models.NewCart()
				h.productServiceMock.Expect(minimock.AnyContext, 123).Return(true, nil)
				h.stocksCheckerMock.Expect(minimock.AnyContext, 123, 1).Return(true, nil)
				h.getCartRepoMock.Expect(minimock.AnyContext, 123).Return(cart, nil)
				h.saveCartRepoMock.Expect(minimock.AnyContext, 123, cart).Return(errorToThrow)
			},
			args: args{
				ctx:   context.Background(),
				user:  123,
				skuId: 123,
				count: 1,
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
			err := helper.service.AddItem(tt.args.ctx, tt.args.user, tt.args.skuId, tt.args.count)
			assert.ErrorIs(t, err, tt.err)
		})
	}
}
