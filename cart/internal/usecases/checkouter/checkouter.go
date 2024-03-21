package checkouter

import (
	"context"
	"route256.ozon.ru/project/cart/internal/models"
)

type repository interface {
	GetCart(ctx context.Context, user int64) (*models.Cart, error)
	SaveCart(ctx context.Context, user int64, cart *models.Cart) error
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
		return 0, err
	}
	items := cart.ListItems(ctx)
	orderId, err := c.orders.CreateOrder(ctx, userId, items)
	if err != nil {
		return 0, err
	}
	cart.Clear(ctx)
	_ = c.repo.SaveCart(ctx, userId, cart)
	return orderId, nil
}
