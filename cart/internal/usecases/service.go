package usecases

import (
	"context"
	"route256.ozon.ru/project/cart/internal/models"
)

type cartadder interface {
	AddItem(ctx context.Context, user int64, skuId int64, count uint16) error
}

type cartdeleter interface {
	DeleteItem(ctx context.Context, user int64, skuId int64) error
}

type cartclearer interface {
	ClearCart(ctx context.Context, user int64) error
}

type cartlister interface {
	ListCartContent(ctx context.Context, user int64) (*models.CartContent, error)
}

type checkouter interface {
	Checkout(ctx context.Context, userId int64) (int64, error)
}

type CartService struct {
	adder      cartadder
	deleter    cartdeleter
	clearer    cartclearer
	lister     cartlister
	checkouter checkouter
}

func NewCartService(adder cartadder, deleter cartdeleter, clearer cartclearer, lister cartlister, checkouter checkouter) *CartService {
	return &CartService{adder: adder, deleter: deleter, clearer: clearer, lister: lister, checkouter: checkouter}
}

func (cs *CartService) AddItem(ctx context.Context, user int64, skuId int64, count uint16) error {
	return cs.adder.AddItem(ctx, user, skuId, count)
}

func (cs *CartService) DeleteItem(ctx context.Context, user int64, skuId int64) error {
	return cs.deleter.DeleteItem(ctx, user, skuId)
}

func (cs *CartService) ClearCart(ctx context.Context, user int64) error {
	return cs.clearer.ClearCart(ctx, user)
}

func (cs *CartService) ListCartContent(ctx context.Context, user int64) (*models.CartContent, error) {
	return cs.lister.ListCartContent(ctx, user)
}

func (cs *CartService) Checkout(ctx context.Context, userId int64) (int64, error) {
	return cs.checkouter.Checkout(ctx, userId)
}
