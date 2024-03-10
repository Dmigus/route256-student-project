package repository

import (
	"context"
	"route256.ozon.ru/project/cart/internal/models"
	"sync"
)

type InMemoryCartRepository struct {
	mu    sync.Mutex
	carts map[int64]*models.Cart
}

func New() *InMemoryCartRepository {
	return &InMemoryCartRepository{
		carts: make(map[int64]*models.Cart),
	}
}

func (c *InMemoryCartRepository) GetCart(_ context.Context, user int64) (*models.Cart, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if cart, exists := c.carts[user]; exists {
		return cart, nil
	}
	newCart := models.NewCart()
	c.carts[user] = newCart
	return newCart, nil
}

func (c *InMemoryCartRepository) SaveCart(_ context.Context, user int64, cart *models.Cart) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.carts[user] = cart
	return nil
}
