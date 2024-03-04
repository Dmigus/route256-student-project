package lister

import (
	"context"
	"route256.ozon.ru/project/cart/internal/service"
)

type CartItem struct {
	SkuId service.SkuId
	Count service.ItemCount
}

type CartToList interface {
	ListItems(ctx context.Context) ([]CartItem, error)
}

type Repository interface {
	CartByUser(ctx context.Context, user service.User) (CartToList, error)
}

type ProductInfo struct {
	name  string
	price service.Price
}

type ProductService interface {
	GetProductsInfo(ctx context.Context, skuIds []service.SkuId) ([]ProductInfo, error)
}

type CartListerService struct {
	repo           Repository
	productService ProductService
}

func (cl *CartListerService) ListCartContent(ctx context.Context, user service.User) (*CartContent, error) {
	cart, err := cl.repo.CartByUser(ctx, user)
	if err != nil {
		return nil, err
	}
	items, err := cart.ListItems(ctx)
	if err != nil {
		return nil, err
	}
	skuIds := extractSkuIds(items)
	productInfos, err := cl.productService.GetProductsInfo(ctx, skuIds)
	if err != nil {
		return nil, err
	}
	return createCartContent(items, productInfos), nil
}

func createCartContent(items []CartItem, prodInfos []ProductInfo) *CartContent {
	content := &CartContent{}
	for i := range items {
		itInfo := ItemInfo{
			SkuId: items[i].SkuId,
			Name:  prodInfos[i].name,
			Count: items[i].Count,
			Price: prodInfos[i].price,
		}
		content.addItem(itInfo)
	}
	return content
}

func extractSkuIds(items []CartItem) []service.SkuId {
	skuIds := make([]service.SkuId, len(items))
	for i, item := range items {
		skuIds[i] = item.SkuId
	}
	return skuIds
}
