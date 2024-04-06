//go:build unit
// +build unit

package productinfogetter

import (
	"github.com/gojuno/minimock/v3"
	"testing"
)

type testHelper struct {
	performMock    *mCallPerformerMockPerform
	prodInfoGetter *ProductInfoGetter
}

func newTestHelper(t *testing.T) testHelper {
	mc := minimock.NewController(t)
	helper := testHelper{}
	performer := NewCallPerformerMock(mc)
	helper.performMock = &performer.PerformMock
	helper.prodInfoGetter = NewProductInfoGetter(performer)
	return helper
}
