package checkouter

import (
	"context"
	"fmt"
	"route256.ozon.ru/project/cart/internal/models"
)

type repository interface {
	GetCart(ctx context.Context, user int64) (*models.Cart, error)
	ClearCartReliable(ctx context.Context, user int64)
}

type orderSystem interface {
	CreateOrder(ctx context.Context, userId int64, items []models.CartItem) (int64, error)
}

type Checkouter struct {
	repo   repository
	orders orderSystem
}

func NewCheckouter(repo repository, orders orderSystem) *Checkouter {
	return &Checkouter{repo: repo, orders: orders}
}

func (c *Checkouter) Checkout(ctx context.Context, userId int64) (int64, error) {
	cart, err := c.repo.GetCart(ctx, userId)
	if err != nil {
		return 0, fmt.Errorf("could not get cart for user %d: %w", userId, err)
	}
	items := cart.ListItemsSorted(ctx)
	orderId, err := c.orders.CreateOrder(ctx, userId, items)
	if err != nil {
		return 0, fmt.Errorf("could not create order for user %d: %w", userId, err)
	}
	c.repo.ClearCartReliable(ctx, userId)
	return orderId, nil
}
