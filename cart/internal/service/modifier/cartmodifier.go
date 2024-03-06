package modifier

import (
	"context"
	"fmt"
	"route256.ozon.ru/project/cart/internal/service"
)

type CartToModify interface {
	Add(ctx context.Context, skuId service.SkuId, count service.ItemCount) error
	Delete(ctx context.Context, skuId service.SkuId) error
	Clear(ctx context.Context) error
}

type Repository interface {
	CartToModifyByUser(ctx context.Context, user service.User) (CartToModify, error)
}

type ProductService interface {
	IsItemPresent(ctx context.Context, skuId service.SkuId) (bool, error)
}

// CartModifierService предназначен для модификации корзин пользователей
type CartModifierService struct {
	repo           Repository
	productService ProductService
}

func New(repo Repository, productService ProductService) *CartModifierService {
	return &CartModifierService{
		repo:           repo,
		productService: productService,
	}
}

func (cs *CartModifierService) AddItem(ctx context.Context, user service.User, skuId service.SkuId, count service.ItemCount) error {
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

func (cs *CartModifierService) DeleteItem(ctx context.Context, user service.User, skuId service.SkuId) error {
	cart, err := cs.repo.CartToModifyByUser(ctx, user)
	if err != nil {
		return err
	}
	return cart.Delete(ctx, skuId)
}

func (cs *CartModifierService) ClearCart(ctx context.Context, user service.User) error {
	cart, err := cs.repo.CartToModifyByUser(ctx, user)
	if err != nil {
		return err
	}
	return cart.Clear(ctx)
}
