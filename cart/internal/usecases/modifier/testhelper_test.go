package modifier

import (
	"github.com/gojuno/minimock/v3"
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
