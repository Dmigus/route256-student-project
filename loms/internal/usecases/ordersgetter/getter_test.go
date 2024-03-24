//go:build unit
// +build unit

package ordersgetter

import (
	"context"
	"fmt"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"route256.ozon.ru/project/loms/internal/models"
	"testing"
)

func TestOrdersGetterPositive(t *testing.T) {
	t.Parallel()
	h := newTestHelper(t)
	order := models.NewOrder(123, 1234)
	order.Status = models.Cancelled
	items := []models.OrderItem{
		{456, 10},
	}
	order.Items = items
	h.orderLoadRepoMock.Expect(minimock.AnyContext, 1234).Return(order, nil)
	returned, err := h.getter.Get(context.Background(), 1234)
	require.NoError(t, err)
	require.NotNil(t, returned)
	assert.Equal(t, models.Cancelled, returned.Status)
	assert.Equal(t, items, returned.Items)
}

func TestOrdersGetterError(t *testing.T) {
	t.Parallel()
	errorToThrow := fmt.Errorf("oops error")
	h := newTestHelper(t)
	h.orderLoadRepoMock.Expect(minimock.AnyContext, 1234).Return(nil, errorToThrow)
	_, err := h.getter.Get(context.Background(), 1234)
	assert.ErrorIs(t, err, errorToThrow)
}
