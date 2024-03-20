//go:build unit
// +build unit

package orderidgenerator

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSequentialGenerator(t *testing.T) {
	t.Parallel()
	gen := NewSequentialGenerator()
	first := gen.NewId()
	require.Equal(t, int64(1), first)
	second := gen.NewId()
	require.Equal(t, int64(2), second)
	assert.Equal(t, int64(2), gen.prevID)
}
