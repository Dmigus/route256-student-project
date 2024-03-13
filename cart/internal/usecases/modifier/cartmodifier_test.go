package modifier

import (
	"context"
	"fmt"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"route256.ozon.ru/project/cart/internal/models"
	"testing"
)

type testHelper struct {
	getCartRepoMock    *mRepositoryMockGetCart
	saveCartRepoMock   *mRepositoryMockSaveCart
	productServiceMock *mProductServiceMockIsItemPresent
	service            *CartModifierService
}

func newTestHelper(t *testing.T) testHelper {
	mc := minimock.NewController(t)
	helper := testHelper{}
	repo := NewRepositoryMock(mc)
	productService := NewProductServiceMock(mc)
	helper.getCartRepoMock = &(repo.GetCartMock)
	helper.saveCartRepoMock = &(repo.SaveCartMock)
	helper.productServiceMock = &(productService.IsItemPresentMock)
	helper.service = New(repo, productService)
	return helper
}

func TestAddItemWithoutErr(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx   context.Context
		user  int64
		skuId int64
		count uint16
	}
	tests := []struct {
		name         string
		args         args
		helperGetter func() testHelper
	}{
		{
			name: "positive",
			args: args{
				ctx:   context.Background(),
				user:  123,
				skuId: 123,
				count: 1,
			},
			helperGetter: func() testHelper {
				h := newTestHelper(t)
				cart := models.NewCart()
				h.getCartRepoMock.Expect(minimock.AnyContext, 123).Return(cart, nil)
				h.productServiceMock.Expect(minimock.AnyContext, 123).Return(true, nil)
				return h
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
			helperGetter: func() testHelper {
				h := newTestHelper(t)
				h.productServiceMock.Expect(minimock.AnyContext, 124).Return(false, nil)
				return h
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			helper := tt.helperGetter()
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
		name         string
		helperGetter func() testHelper
		args         args
		err          error
	}{
		{
			name: "error checking good presence",
			helperGetter: func() testHelper {
				h := newTestHelper(t)
				h.productServiceMock.Expect(minimock.AnyContext, 123).Return(false, errorToThrow)
				return h
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
			name: "error getting user cart",
			helperGetter: func() testHelper {
				h := newTestHelper(t)
				h.productServiceMock.Expect(minimock.AnyContext, 123).Return(true, nil)
				h.getCartRepoMock.Expect(minimock.AnyContext, 123).Return(nil, errorToThrow)
				return h
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
			helperGetter: func() testHelper {
				h := newTestHelper(t)
				cart := models.NewCart()
				h.productServiceMock.Expect(minimock.AnyContext, 123).Return(true, nil)
				h.getCartRepoMock.Expect(minimock.AnyContext, 123).Return(cart, nil)
				h.saveCartRepoMock.Expect(minimock.AnyContext, 123, cart).Return(errorToThrow)
				return h
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
			helper := tt.helperGetter()
			err := helper.service.AddItem(tt.args.ctx, tt.args.user, tt.args.skuId, tt.args.count)
			assert.ErrorIs(t, err, tt.err)
		})
	}
}

func TestCartModifierService_DeleteItem(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)
	errorToThrow := fmt.Errorf("oops error")
	type fields struct {
		repo           repository
		productService productService
	}
	type args struct {
		ctx   context.Context
		user  int64
		skuId int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "positive",
			fields: fields{
				repo: NewRepositoryMock(mc).
					GetCartMock.
					When(minimock.AnyContext, 123).
					Then(models.NewCart(), nil).
					SaveCartMock.
					Return(nil),
			},
			args: args{
				context.Background(),
				123,
				123,
			},
			wantErr: nil,
		},
		{
			name: "error getting user cart",
			fields: fields{
				repo: NewRepositoryMock(mc).
					GetCartMock.
					When(minimock.AnyContext, 123).
					Then(nil, errorToThrow),
			},
			args: args{
				context.Background(),
				123,
				123,
			},
			wantErr: errorToThrow,
		},
		{
			name: "error saving user cart",
			fields: fields{
				repo: NewRepositoryMock(mc).
					GetCartMock.
					When(minimock.AnyContext, 123).
					Then(models.NewCart(), nil).
					SaveCartMock.
					Return(errorToThrow),
			},
			args: args{
				context.Background(),
				123,
				123,
			},
			wantErr: errorToThrow,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cs := &CartModifierService{
				repo:           tt.fields.repo,
				productService: tt.fields.productService,
			}
			err := cs.DeleteItem(tt.args.ctx, tt.args.user, tt.args.skuId)
			assert.ErrorIs(t, err, tt.wantErr, fmt.Errorf("CartModifierService.DeleteItem() error = %v, wantErr %v", err, tt.wantErr))
		})
	}
}

func TestCartModifierService_ClearCart(t *testing.T) {
	mc := minimock.NewController(t)
	errorToThrow := fmt.Errorf("oops error")
	type fields struct {
		repo           repository
		productService productService
	}
	type args struct {
		ctx  context.Context
		user int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "positive",
			fields: fields{
				repo: NewRepositoryMock(mc).
					GetCartMock.
					When(minimock.AnyContext, 123).
					Then(models.NewCart(), nil).
					SaveCartMock.
					Return(nil),
			},
			args: args{
				context.Background(),
				123,
			},
			wantErr: nil,
		},
		{
			name: "error getting user cart",
			fields: fields{
				repo: NewRepositoryMock(mc).
					GetCartMock.
					When(minimock.AnyContext, 123).
					Then(nil, errorToThrow),
			},
			args: args{
				context.Background(),
				123,
			},
			wantErr: errorToThrow,
		},
		{
			name: "error saving user cart",
			fields: fields{
				repo: NewRepositoryMock(mc).
					GetCartMock.
					When(minimock.AnyContext, 123).
					Then(models.NewCart(), nil).
					SaveCartMock.
					Return(errorToThrow),
			},
			args: args{
				context.Background(),
				123,
			},
			wantErr: errorToThrow,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cs := &CartModifierService{
				repo:           tt.fields.repo,
				productService: tt.fields.productService,
			}
			err := cs.ClearCart(tt.args.ctx, tt.args.user)
			assert.ErrorIs(t, err, tt.wantErr, fmt.Errorf("CartModifierService.ClearCart() error = %v, wantErr %v", err, tt.wantErr))
		})
	}
}
