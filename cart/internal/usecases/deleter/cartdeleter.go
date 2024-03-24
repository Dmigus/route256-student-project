package deleter

import (
	"context"
	"fmt"
	"route256.ozon.ru/project/cart/internal/models"
)

type repository interface {
	GetCart(ctx context.Context, user int64) (*models.Cart, error)
	SaveCart(ctx context.Context, user int64, cart *models.Cart) error
}

// CartDeleter предназначен для удаления товаров из корзин пользователей
type CartDeleter struct {
	repo repository
}

func NewCartDeleter(repo repository) *CartDeleter {
	return &CartDeleter{repo: repo}
}
func (c *CartDeleter) DeleteItem(ctx context.Context, user int64, skuId int64) error {
	cart, err := c.repo.GetCart(ctx, user)
	if err != nil {
		return fmt.Errorf("could not get cart for user %d: %w", user, err)
	}
	cart.Delete(ctx, skuId)
	err = c.repo.SaveCart(ctx, user, cart)
	if err != nil {
		return fmt.Errorf("could not get save for user %d: %w", user, err)
	}
	return nil
}
