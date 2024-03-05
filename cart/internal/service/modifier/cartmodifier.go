package modifier

import (
	"context"
	"errors"
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

var ErrItemNotExists = errors.New("item is not exist")

func (cs *CartModifierService) AddItem(ctx context.Context, user service.User, skuId service.SkuId, count service.ItemCount) error {
	if exists, err := cs.productService.IsItemPresent(ctx, skuId); err != nil {
		return fmt.Errorf("could not check item %d presence: %w", skuId, err)
	} else if !exists {
		return ErrItemNotExists
	}
	if cart, err := cs.repo.CartToModifyByUser(ctx, user); err != nil {
		return err
	} else {
		return cart.Add(ctx, skuId, count)
	}
}

func (cs *CartModifierService) DeleteItem(ctx context.Context, user service.User, skuId service.SkuId) error {
	if cart, err := cs.repo.CartToModifyByUser(ctx, user); err != nil {
		return err
	} else {
		return cart.Delete(ctx, skuId)
	}
}

func (cs *CartModifierService) ClearCart(ctx context.Context, user service.User) error {
	if cart, err := cs.repo.CartToModifyByUser(ctx, user); err != nil {
		return err
	} else {
		return cart.Clear(ctx)
	}
}
