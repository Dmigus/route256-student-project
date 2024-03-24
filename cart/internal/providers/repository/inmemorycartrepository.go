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

func (i *InMemoryCartRepository) GetCart(_ context.Context, user int64) (*models.Cart, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if cart, exists := i.carts[user]; exists {
		return cart, nil
	}
	newCart := models.NewCart()
	i.carts[user] = newCart
	return newCart, nil
}

func (i *InMemoryCartRepository) SaveCart(_ context.Context, user int64, cart *models.Cart) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.carts[user] = cart
	return nil
}

// ClearCartReliable "надёжно" очищает корзину для пользователя. Сразу после возрата из метода корзина пользователя гарантированно пуста.
func (i *InMemoryCartRepository) ClearCartReliable(ctx context.Context, user int64) {
	cart, _ := i.GetCart(ctx, user)
	cart.Clear(ctx)
	i.mu.Lock()
	defer i.mu.Unlock()
	i.carts[user] = cart
}
