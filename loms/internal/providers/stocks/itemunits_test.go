//go:build unit
// +build unit

package stocks

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestItemUnits_cancelReserve(t *testing.T) {
	t.Parallel()
	iu := NewItemUnits(10, 5)
	iu.cancelReserve(3)
	assert.Equal(t, uint64(8), iu.getNumOfAvailable())
}

func TestItemUnits_removeReserved(t *testing.T) {
	t.Parallel()
	iu := NewItemUnits(10, 5)
	availableBefore := iu.getNumOfAvailable()
	iu.removeReserved(3)
	assert.Equal(t, availableBefore, iu.getNumOfAvailable())
	assert.Equal(t, uint64(7), iu.total)
	assert.Equal(t, uint64(2), iu.reserved)
}

func TestItemUnits_reserve(t *testing.T) {
	t.Parallel()
	iu := NewItemUnits(10, 5)
	err := iu.reserve(3)
	require.NoError(t, err)
	assert.Equal(t, uint64(2), iu.getNumOfAvailable())
	err = iu.reserve(3)
	require.ErrorIs(t, err, ErrNotEnoughItems)
}
