//go:build unit
// +build unit

package stocksinfogetter

import (
	"context"
	"fmt"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPositive(t *testing.T) {
	t.Parallel()
	helper := newTestHelper(t)
	helper.stockRepoMock.Expect(minimock.AnyContext, 1234).Return(10, nil)
	returned, err := helper.getter.GetNumOfAvailable(context.Background(), 1234)
	require.NoError(t, err)
	assert.Equal(t, uint64(10), returned)
}

func TestError(t *testing.T) {
	t.Parallel()
	someErr := fmt.Errorf("some error")
	helper := newTestHelper(t)
	helper.stockRepoMock.Expect(minimock.AnyContext, 1234).Return(0, someErr)
	_, err := helper.getter.GetNumOfAvailable(context.Background(), 1234)
	assert.ErrorIs(t, err, someErr)
}
