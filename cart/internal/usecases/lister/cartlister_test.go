package lister

import (
	"context"
	"fmt"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sort"
	"testing"

	"route256.ozon.ru/project/cart/internal/models"
)

type testHelper struct {
	repoMock           *mRepositoryMockGetCart
	productServiceMock *mProductServiceMockGetProductsInfo
	service            *CartListerService
}

func newTestHelper(t *testing.T) testHelper {
	mc := minimock.NewController(t)
	helper := testHelper{}
	repo := NewRepositoryMock(mc)
	productService := NewProductServiceMock(mc)
	helper.repoMock = &(repo.GetCartMock)
	helper.productServiceMock = &(productService.GetProductsInfoMock)
	helper.service = New(repo, productService)
	return helper
}

func TestCartListerService_ListCartContentErrors(t *testing.T) {
	t.Parallel()
	errorToThrow := fmt.Errorf("oops error")
	type args struct {
		ctx  context.Context
		user int64
	}
	tests := []struct {
		name         string
		helperGetter func() testHelper
		args         args
		wantErr      error
	}{
		{
			name: "error getting cart",
			helperGetter: func() testHelper {
				helper1 := newTestHelper(t)
				helper1.repoMock.Expect(minimock.AnyContext, 123).Return(nil, errorToThrow)
				return helper1
			},
			args:    args{context.Background(), 123},
			wantErr: errorToThrow,
		},
		{
			name: "error getting products info",
			helperGetter: func() testHelper {
				helper2 := newTestHelper(t)
				helper2.repoMock.Expect(minimock.AnyContext, 123).Return(models.NewCart(), nil)
				helper2.productServiceMock.Return(nil, errorToThrow)
				return helper2
			},
			args:    args{context.Background(), 123},
			wantErr: errorToThrow,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			helper := tt.helperGetter()
			_, err := helper.service.ListCartContent(tt.args.ctx, tt.args.user)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestCartListerService_ListCartContentPositive(t *testing.T) {
	userId := int64(12345)
	cart := models.NewCart()
	cart.Add(context.Background(), 123, 2)
	cart.Add(context.Background(), 456, 10)
	mc := minimock.NewController(t)
	repo := NewRepositoryMock(mc)
	repo.GetCartMock.Expect(minimock.AnyContext, userId).Return(cart, nil)

	prodService := NewProductServiceMock(mc)
	prodInfos := []models.ProductInfo{
		{Name: "item123", Price: 100},
		{Name: "item456", Price: 50},
	}
	prodService.GetProductsInfoMock.Set(func(_ context.Context, skuIds []int64) (pa1 []models.ProductInfo, err error) {
		require.Len(t, skuIds, 2, "passed skuIds have wrong len")
		if skuIds[0] > skuIds[1] {
			prodInfos[0], prodInfos[1] = prodInfos[1], prodInfos[0]
		}
		return prodInfos, nil
	})
	service := New(repo, prodService)

	want := models.NewCartContent()
	want.Add(models.CartItemInfo{
		CartItem:    models.CartItem{SkuId: 123, Count: 2},
		ProductInfo: models.ProductInfo{Name: "item123", Price: 100},
	})
	want.Add(models.CartItemInfo{
		CartItem:    models.CartItem{SkuId: 456, Count: 10},
		ProductInfo: models.ProductInfo{Name: "item456", Price: 50},
	})

	got, err := service.ListCartContent(context.Background(), userId)
	sort.Slice(got.GetItems(), func(i, j int) bool {
		return got.GetItems()[i].CartItem.SkuId < got.GetItems()[j].CartItem.SkuId
	})
	require.NoError(t, err, "must be nil err in positive scenario")
	require.Equal(t, want, got)
}

func Test_createCartContent(t *testing.T) {
	t.Parallel()
	items := []models.CartItem{
		{SkuId: 123, Count: 1},
		{SkuId: 456, Count: 10},
	}
	prodInfos := []models.ProductInfo{
		{Name: "item123", Price: 100},
		{Name: "item456", Price: 200},
	}

	want := &models.CartContent{}
	want.Add(models.CartItemInfo{
		CartItem:    models.CartItem{SkuId: 123, Count: 1},
		ProductInfo: models.ProductInfo{Name: "item123", Price: 100},
	})
	want.Add(models.CartItemInfo{
		CartItem:    models.CartItem{SkuId: 456, Count: 10},
		ProductInfo: models.ProductInfo{Name: "item456", Price: 200},
	})

	got := createCartContent(items, prodInfos)
	assert.Equal(t, want, got, fmt.Sprintf("createCartContent() = %v, want %v", got, want))
}

func Test_extractSkuIds(t *testing.T) {
	t.Parallel()
	type args struct {
		items []models.CartItem
	}
	tests := []struct {
		name string
		args args
		want []int64
	}{
		{
			name: "three elems",
			args: args{
				[]models.CartItem{
					{SkuId: 1, Count: 1},
					{SkuId: 2, Count: 5},
					{SkuId: 3, Count: 6},
				},
			},
			want: []int64{1, 2, 3},
		},
		{
			name: "empty",
			args: args{
				[]models.CartItem{},
			},
			want: []int64{},
		},
		{
			name: "nil",
			args: args{nil},
			want: []int64{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := extractSkuIds(tt.args.items)
			assert.Equal(t, tt.want, got, fmt.Sprintf("extractSkuIds() = %v, want %v", got, tt.want))
		})
	}
}
