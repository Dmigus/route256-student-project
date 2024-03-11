package repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"route256.ozon.ru/project/cart/internal/models"
	"testing"
)

func TestInMemoryCartRepository_GetCart(t *testing.T) {
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
	userId := int64(123)
	err := repo.SaveCart(context.Background(), userId, newCart)
	require.NoError(t, err, "saving cart failed with error")
	returnedCart, err := repo.GetCart(context.Background(), userId)
	require.NoError(t, err, "getting cart failed with error")
	assert.True(t, returnedCart == newCart, "returned card is not same as saved for this user")
	otherUserId := int64(456)
	returnedOtherCard, err := repo.GetCart(context.Background(), otherUserId)
	require.NoError(t, err, "getting cart failed with error")
	assert.False(t, returnedOtherCard == newCart, "the same cart returned for different users")
}
