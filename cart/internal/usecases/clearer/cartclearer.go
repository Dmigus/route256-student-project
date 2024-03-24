package clearer

import (
	"context"
	"fmt"
	"route256.ozon.ru/project/cart/internal/models"
)

type repository interface {
	GetCart(ctx context.Context, user int64) (*models.Cart, error)
	SaveCart(ctx context.Context, user int64, cart *models.Cart) error
}

// CartClearer предназначен для очищения корзин пользователей
type CartClearer struct {
	repo repository
}

func NewCartClearer(repo repository) *CartClearer {
	return &CartClearer{repo: repo}
}

func (c *CartClearer) ClearCart(ctx context.Context, user int64) error {
	cart, err := c.repo.GetCart(ctx, user)
	if err != nil {
		return fmt.Errorf("could not get cart for user %d: %w", user, err)
	}
	cart.Clear(ctx)
	err = c.repo.SaveCart(ctx, user, cart)
	if err != nil {
		return fmt.Errorf("could not save cart for user %d: %w", user, err)
	}
	return nil
}
