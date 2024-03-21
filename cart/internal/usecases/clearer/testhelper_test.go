//go:build unit
// +build unit

package clearer

import (
	"github.com/gojuno/minimock/v3"
	"testing"
)

type testHelper struct {
	getCartRepoMock  *mRepositoryMockGetCart
	saveCartRepoMock *mRepositoryMockSaveCart
	service          *CartClearer
}

func newTestHelper(t *testing.T) testHelper {
	mc := minimock.NewController(t)
	helper := testHelper{}
	repo := NewRepositoryMock(mc)
	helper.getCartRepoMock = &(repo.GetCartMock)
	helper.saveCartRepoMock = &(repo.SaveCartMock)
	helper.service = NewCartClearer(repo)
	return helper
}
