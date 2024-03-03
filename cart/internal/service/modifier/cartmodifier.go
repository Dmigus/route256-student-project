package modifier

import (
	"context"
	"errors"
	"fmt"
)

type (
	User  = int64
	SkuId = int64
)

type CartToModify interface {
	Add(ctx context.Context, skuId SkuId, count uint16) error
	Delete(ctx context.Context, skuId SkuId) error
	Clear(ctx context.Context) error
}

type Repository interface {
	CartByUser(ctx context.Context, user User) (CartToModify, error)
}

type ProductService interface {
	IsItemPresent(ctx context.Context, skuId SkuId) (bool, error)
}

// Сервис, предназначенный для модификации корзин пользователей
type CartModifierService struct {
	repo           Repository
	productService ProductService
}

var ErrItemNotExists = errors.New("item is not exist")

func (cs *CartModifierService) AddItem(ctx context.Context, user User, skuId SkuId, count uint16) error {
	if exists, err := cs.productService.IsItemPresent(ctx, skuId); err != nil {
		return fmt.Errorf("could not check item %d presence: %w", skuId, err)
	} else if !exists {
		return ErrItemNotExists
	}
	if cart, err := cs.repo.CartByUser(ctx, user); err != nil {
		return err
	} else {
		return cart.Add(ctx, skuId, count)
	}
}

func (cs *CartModifierService) DeleteItem(ctx context.Context, user User, skuId SkuId) error {
	if cart, err := cs.repo.CartByUser(ctx, user); err != nil {
		return err
	} else {
		return cart.Delete(ctx, skuId)
	}
}

func (cs *CartModifierService) ClearCart(ctx context.Context, user User) error {
	if cart, err := cs.repo.CartByUser(ctx, user); err != nil {
		return err
	} else {
		return cart.Clear(ctx)
	}
}
