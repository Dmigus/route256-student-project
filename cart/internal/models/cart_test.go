//go:build unit
// +build unit

package models

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCart(t *testing.T) {
	t.Parallel()
	cart := NewCart()
	ctx := context.Background()
	cart.Add(ctx, 123, 10)
	cart.Add(ctx, 123, 15)
	items := cart.ListItems(ctx)
	require.Len(t, items, 1)
	assert.Equal(t, []CartItem{{123, 25}}, items)

	cart.Delete(ctx, 123)
	items = cart.ListItems(ctx)
	require.Len(t, items, 0)

	cart.Add(ctx, 1234, 15)
	cart.Add(ctx, 123, 10)
	items = cart.ListItemsSorted(ctx)
	assert.Equal(t, []CartItem{{123, 10}, {1234, 15}}, items)
	cart.Clear(ctx)
	items = cart.ListItems(ctx)
	assert.Len(t, items, 0)
}
