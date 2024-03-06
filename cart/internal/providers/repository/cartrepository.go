package repository

import (
	"context"
	"route256.ozon.ru/project/cart/internal/usecases"
	"route256.ozon.ru/project/cart/internal/usecases/lister"
	"route256.ozon.ru/project/cart/internal/usecases/modifier"
	"sync"
)

type Cart interface {
	lister.CartToList
	modifier.CartToModify
}

type cartCreator interface {
	Create(ctx context.Context) (Cart, error)
}

type CartRepository struct {
	cartCreator cartCreator
	carts       map[usecases.User]Cart
	mu          sync.Mutex
}

func New(cartCreator cartCreator) *CartRepository {
	return &CartRepository{
		cartCreator: cartCreator,
		carts:       make(map[usecases.User]Cart),
	}
}

func (c *CartRepository) getCartByUser(ctx context.Context, user usecases.User) (Cart, error) {
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

func (c *CartRepository) CartToListByUser(ctx context.Context, user usecases.User) (lister.CartToList, error) {
	return c.getCartByUser(ctx, user)
}

func (c *CartRepository) CartToModifyByUser(ctx context.Context, user usecases.User) (modifier.CartToModify, error) {
	return c.getCartByUser(ctx, user)
}
