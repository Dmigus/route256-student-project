//go:build unit
// +build unit

package lister

import (
	"github.com/gojuno/minimock/v3"
	"testing"
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
