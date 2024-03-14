package repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"route256.ozon.ru/project/cart/internal/models"
	"sort"
	"testing"
)

func TestInMemoryCartRepository_GetNewCarts(t *testing.T) {
	t.Parallel()
	repo := New()
	userId := int64(123)
	newCart, err := repo.GetCart(context.Background(), userId)
	require.NoError(t, err, "cart creation failed")
	assert.Len(t, newCart.ListItems(context.Background()), 0, "created cart is not empty")
	sameCart, err := repo.GetCart(context.Background(), userId)
	require.NoError(t, err, "created cart is not empty")
	assert.True(t, newCart == sameCart, "returned different carts on same user")
	otherUser := int64(456)
	otherUserCart, err := repo.GetCart(context.Background(), otherUser)
	require.NoError(t, err, "created cart is not empty")
	assert.False(t, newCart == otherUserCart, "the same cart returned for different users")
}

func TestInMemoryCartRepository_GetExistingCart(t *testing.T) {
	t.Parallel()
	repo := New()
	newCart := models.NewCart()
	newCart.Add(context.Background(), 456, 10)
	newCart.Add(context.Background(), 789, 15)
	userId := int64(123)
	repo.carts[userId] = newCart
	returnedCart, err := repo.GetCart(context.Background(), userId)
	require.NoError(t, err, "getting cart failed with error")
	returnedItems := returnedCart.ListItems(context.Background())
	expectedItems := []models.CartItem{
		{456, 10},
		{789, 15},
	}
	sort.Slice(returnedItems, func(i, j int) bool {
		return returnedItems[i].SkuId < returnedItems[j].SkuId
	})
	assert.Equal(t, expectedItems, returnedItems)
}

func TestInMemoryCartRepository_SaveCart(t *testing.T) {
	t.Parallel()
	repo := New()
	newCart := models.NewCart()
	userId := int64(123)
	err := repo.SaveCart(context.Background(), userId, newCart)
	require.NoError(t, err, "saving cart failed with error")
	assert.True(t, newCart == repo.carts[userId], "saving failed")
}

func TestInMemoryCartRepository_SaveAndGet(t *testing.T) {
	t.Parallel()
	repo := New()
	newCart := models.NewCart()
	newCart.Add(context.Background(), 456, 10)
	newCart.Add(context.Background(), 789, 15)
	userId := int64(123)
	err := repo.SaveCart(context.Background(), userId, newCart)
	require.NoError(t, err, "saving cart failed with error")
	returnedCart, err := repo.GetCart(context.Background(), userId)
	require.NoError(t, err, "getting cart failed with error")
	assert.Equal(t, newCart, returnedCart, "returned card is not same as saved for this user")

	otherUserId := int64(456)
	returnedOtherCard, err := repo.GetCart(context.Background(), otherUserId)
	require.NoError(t, err, "getting cart failed with error")
	assert.NotEqual(t, newCart, returnedOtherCard, "the same cart returned for different users")
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
