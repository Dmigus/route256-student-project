//go:build unit
// +build unit

package adder

import (
	"github.com/gojuno/minimock/v3"
	"testing"
)

type testHelper struct {
	getCartRepoMock    *mRepositoryMockGetCart
	saveCartRepoMock   *mRepositoryMockSaveCart
	productServiceMock *mProductServiceMockIsItemPresent
	stocksCheckerMock  *mStocksCheckerMockIsItemAvailable
	service            *CartAdderService
}

func newTestHelper(t *testing.T) testHelper {
	mc := minimock.NewController(t)
	helper := testHelper{}
	repo := NewRepositoryMock(mc)
	productService := NewProductServiceMock(mc)
	stocks := NewStocksCheckerMock(mc)
	helper.getCartRepoMock = &(repo.GetCartMock)
	helper.saveCartRepoMock = &(repo.SaveCartMock)
	helper.productServiceMock = &(productService.IsItemPresentMock)
	helper.stocksCheckerMock = &(stocks.IsItemAvailableMock)
	helper.service = New(repo, productService, stocks)
	return helper
}
