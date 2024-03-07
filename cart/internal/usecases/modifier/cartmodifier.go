package modifier

import (
	"context"
	"fmt"
	"route256.ozon.ru/project/cart/internal/models"
)

type repository interface {
	GetCart(ctx context.Context, user models.UserId) (*models.Cart, error)
	SaveCart(ctx context.Context, user models.UserId, cart *models.Cart) error
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
	cart, err := cs.repo.GetCart(ctx, user)
	if err != nil {
		return err
	}
	cart.Add(ctx, skuId, count)
	return cs.repo.SaveCart(ctx, user, cart)
}

func (cs *CartModifierService) DeleteItem(ctx context.Context, user models.UserId, skuId models.SkuId) error {
	cart, err := cs.repo.GetCart(ctx, user)
	if err != nil {
		return err
	}
	cart.Delete(ctx, skuId)
	return cs.repo.SaveCart(ctx, user, cart)
}

func (cs *CartModifierService) ClearCart(ctx context.Context, user models.UserId) error {
	cart, err := cs.repo.GetCart(ctx, user)
	if err != nil {
		return err
	}
	cart.Clear(ctx)
	return cs.repo.SaveCart(ctx, user, cart)
}
