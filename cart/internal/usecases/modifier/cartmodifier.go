package modifier

import (
	"context"
	"fmt"
	"route256.ozon.ru/project/cart/internal/usecases"
)

type CartToModify interface {
	Add(ctx context.Context, skuId usecases.SkuId, count usecases.ItemCount) error
	Delete(ctx context.Context, skuId usecases.SkuId) error
	Clear(ctx context.Context) error
}

type repository interface {
	CartToModifyByUser(ctx context.Context, user usecases.User) (CartToModify, error)
}

type productService interface {
	IsItemPresent(ctx context.Context, skuId usecases.SkuId) (bool, error)
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

func (cs *CartModifierService) AddItem(ctx context.Context, user usecases.User, skuId usecases.SkuId, count usecases.ItemCount) error {
	exists, err := cs.productService.IsItemPresent(ctx, skuId)
	if err != nil {
		return fmt.Errorf("could not check item %d presence: %w", skuId, err)
	}
	if !exists {
		return nil
	}
	cart, err := cs.repo.CartToModifyByUser(ctx, user)
	if err != nil {
		return err
	}
	return cart.Add(ctx, skuId, count)
}

func (cs *CartModifierService) DeleteItem(ctx context.Context, user usecases.User, skuId usecases.SkuId) error {
	cart, err := cs.repo.CartToModifyByUser(ctx, user)
	if err != nil {
		return err
	}
	return cart.Delete(ctx, skuId)
}

func (cs *CartModifierService) ClearCart(ctx context.Context, user usecases.User) error {
	cart, err := cs.repo.CartToModifyByUser(ctx, user)
	if err != nil {
		return err
	}
	return cart.Clear(ctx)
}
