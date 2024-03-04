package cartrepository

import (
	"context"
	"route256.ozon.ru/project/cart/internal/service"
	"route256.ozon.ru/project/cart/internal/service/lister"
	"route256.ozon.ru/project/cart/internal/service/modifier"
	"sync"
)

type Cart interface {
	lister.CartToList
	modifier.CartToModify
}

type CartFabric interface {
	Create(ctx context.Context) (Cart, error)
}

type CartRepository struct {
	cartFabric CartFabric
	carts      map[service.User]Cart
	mu         sync.Mutex
}

func New(cartFabric CartFabric) *CartRepository {
	return &CartRepository{
		cartFabric: cartFabric,
		carts:      make(map[service.User]Cart),
	}
}

func (c *CartRepository) CartByUser(ctx context.Context, user service.User) (Cart, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if cart, exists := c.carts[user]; exists {
		return cart, nil
	}
	newCart, err := c.cartFabric.Create(ctx)
	if err != nil {
		return nil, err
	}
	c.carts[user] = newCart
	return newCart, nil
}
