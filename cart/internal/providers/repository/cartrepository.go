package repository

import (
	"context"
	"route256.ozon.ru/project/cart/internal/models"
	"sync"
)

type CartRepository struct {
	carts map[models.UserId]*models.InMemoryCart
	mu    sync.Mutex
}

func New() *CartRepository {
	return &CartRepository{
		carts: make(map[models.UserId]*models.InMemoryCart),
	}
}

func (c *CartRepository) CartByUser(_ context.Context, user models.UserId) (*models.InMemoryCart, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if cart, exists := c.carts[user]; exists {
		return cart, nil
	}
	newCart := models.NewInMemoryCart()
	c.carts[user] = newCart
	return newCart, nil
}
