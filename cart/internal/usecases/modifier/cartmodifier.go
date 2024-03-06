package modifier

import (
	"context"
	"fmt"
	"route256.ozon.ru/project/cart/internal/models"
)

type repository interface {
	CartByUser(ctx context.Context, user models.UserId) (*models.InMemoryCart, error)
}

type productService interface {
	IsItemPresent(ctx context.Context, skuId models.SkuId) (bool, error)
}

// CartModifierService предназначен для модификации корзин пользователей
type CartModifierService struct {
	repo           repository
	productService productService
}

func New(repo repository, productService productService) *CartModifierService {
	return &CartModifierService{
		repo:           repo,
		productService: productService,
	}
}

func (cs *CartModifierService) AddItem(ctx context.Context, user models.UserId, skuId models.SkuId, count models.ItemCount) error {
	exists, err := cs.productService.IsItemPresent(ctx, skuId)
	if err != nil {
		return fmt.Errorf("could not check item %d presence: %w", skuId, err)
	}
	if !exists {
		return nil
	}
	cart, err := cs.repo.CartByUser(ctx, user)
	if err != nil {
		return err
	}
	return cart.Add(ctx, skuId, count)
}

func (cs *CartModifierService) DeleteItem(ctx context.Context, user models.UserId, skuId models.SkuId) error {
	cart, err := cs.repo.CartByUser(ctx, user)
	if err != nil {
		return err
	}
	return cart.Delete(ctx, skuId)
}

func (cs *CartModifierService) ClearCart(ctx context.Context, user models.UserId) error {
	cart, err := cs.repo.CartByUser(ctx, user)
	if err != nil {
		return err
	}
	return cart.Clear(ctx)
}
