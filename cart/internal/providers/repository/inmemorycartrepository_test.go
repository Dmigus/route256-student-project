//go:build unit
// +build unit

package repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"route256.ozon.ru/project/cart/internal/models"
	"testing"
)

func TestInMemoryCartRepository_GetExistingCart(t *testing.T) {
	t.Parallel()
	repo := New()
	ctx := context.Background()
	newCart := models.NewCart()
	newCart.Add(ctx, 456, 10)
	newCart.Add(ctx, 789, 15)
	userId := int64(123)
	repo.carts[userId] = newCart
	returnedCart, err := repo.GetCart(ctx, userId)
	require.NoError(t, err, "getting cart failed with error")
	returnedItems := returnedCart.ListItemsSorted(context.Background())
	expectedItems := []models.CartItem{
		{456, 10},
		{789, 15},
	}
	assert.Equal(t, expectedItems, returnedItems)
}

func TestInMemoryCartRepository_SaveAndGet(t *testing.T) {
	t.Parallel()
	repo := New()
	newCart := models.NewCart()
	newCart.Add(context.Background(), 456, 10)
	newCart.Add(context.Background(), 789, 15)
	userId := int64(123)
	ctx := context.Background()
	err := repo.SaveCart(ctx, userId, newCart)
	require.NoError(t, err, "saving cart failed with error")
	returnedCart, err := repo.GetCart(context.Background(), userId)
	require.NoError(t, err, "getting cart failed with error")
	assert.Equal(t, newCart.ListItemsSorted(ctx), returnedCart.ListItemsSorted(ctx), "returned card is not same as saved for this user")

	otherUserId := int64(456)
	returnedOtherCard, err := repo.GetCart(context.Background(), otherUserId)
	require.NoError(t, err, "getting cart failed with error")
	assert.NotEqual(t, newCart.ListItemsSorted(ctx), returnedOtherCard.ListItemsSorted(ctx), "the same cart returned for different users")
}

func BenchmarkGetNewCarts(b *testing.B) {
	repo := New()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		randUserId := rand.Int63()
		b.StartTimer()
		_, _ = repo.GetCart(ctx, randUserId)
	}
}

func BenchmarkGetSameCart(b *testing.B) {
	repo := New()
	ctx := context.Background()
	userId := int64(123)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.GetCart(ctx, userId)
	}
}

func BenchmarkSaveNewCarts(b *testing.B) {
	repo := New()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		randUserId := rand.Int63()
		newCart := models.NewCart()
		b.StartTimer()
		_ = repo.SaveCart(ctx, randUserId, newCart)
	}
}
