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

type CartCreator interface {
	Create(ctx context.Context) (Cart, error)
}

type CartRepository struct {
	cartCreator CartCreator
	carts       map[service.User]Cart
	mu          sync.Mutex
}

func New(cartCreator CartCreator) *CartRepository {
	return &CartRepository{
		cartCreator: cartCreator,
		carts:       make(map[service.User]Cart),
	}
}

func (c *CartRepository) getCartByUser(ctx context.Context, user service.User) (Cart, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if cart, exists := c.carts[user]; exists {
		return cart, nil
	}
	newCart, err := c.cartCreator.Create(ctx)
	if err != nil {
		return nil, err
	}
	c.carts[user] = newCart
	return newCart, nil
}

func (c *CartRepository) CartToListByUser(ctx context.Context, user service.User) (lister.CartToList, error) {
	return c.getCartByUser(ctx, user)
}

func (c *CartRepository) CartToModifyByUser(ctx context.Context, user service.User) (modifier.CartToModify, error) {
	return c.getCartByUser(ctx, user)
}
