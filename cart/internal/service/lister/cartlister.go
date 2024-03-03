package lister

import (
	"context"
	"route256.ozon.ru/project/cart/internal/service/modifier"
)

type CartToList interface {
	Range(ctx context.Context, f func(ctx context.Context, skuId modifier.SkuId, count uint16) bool)
}

type Repository interface {
	CartByUser(ctx context.Context, user modifier.User) (CartToList, error)
}

type productInfo struct {
	name  string
	price uint32
}

type ProductService interface {
	GetProductInfo(ctx context.Context, skuId modifier.SkuId) (productInfo, error)
}

type CartListerService struct {
	repo           Repository
	productService ProductService
}

func (cl *CartListerService) ListCartContent(ctx context.Context, user modifier.User) (CartContent, error) {
	cart, err := cl.repo.CartByUser(ctx, user)
	if err != nil {
		return CartContent{}, err
	}
	return cl.compCartContent(ctx, cart)
}

func (cl *CartListerService) compCartContent(ctx context.Context, cart CartToList) (CartContent, error) {
	content := CartContent{}
	var errGlob error
	cart.Range(ctx, func(ctx context.Context, skuId modifier.SkuId, count uint16) bool {
		prodInfo, err := cl.productService.GetProductInfo(ctx, skuId)
		if err != nil {
			errGlob = err
			return false
		}
		itInfo := ItemInfo{
			SkuId: skuId,
			Name:  prodInfo.name,
			Count: count,
			Price: prodInfo.price,
		}
		content.items = append(content.items, itInfo)
		content.totalPrice += uint32(count) * prodInfo.price
		return true
	})
	if errGlob != nil {
		return CartContent{}, errGlob
	}
	return content, nil
}
