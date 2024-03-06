package lister

import (
	"context"
	"route256.ozon.ru/project/cart/internal/usecases"
)

type CartItem struct {
	SkuId usecases.SkuId
	Count usecases.ItemCount
}

type CartToList interface {
	ListItems(ctx context.Context) ([]CartItem, error)
}

type repository interface {
	CartToListByUser(ctx context.Context, user usecases.User) (CartToList, error)
}

type ProductInfo struct {
	Name  string
	Price usecases.Price
}

type productService interface {
	GetProductsInfo(ctx context.Context, skuIds []usecases.SkuId) ([]ProductInfo, error)
}

type CartListerService struct {
	repo           repository
	productService productService
}

func New(repo repository, productService productService) *CartListerService {
	return &CartListerService{repo: repo, productService: productService}
}

func (cl *CartListerService) ListCartContent(ctx context.Context, user usecases.User) (*CartContent, error) {
	cart, err := cl.repo.CartToListByUser(ctx, user)
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
			Name:  prodInfos[i].Name,
			Count: items[i].Count,
			Price: prodInfos[i].Price,
		}
		content.addItem(itInfo)
	}
	return content
}

func extractSkuIds(items []CartItem) []usecases.SkuId {
	skuIds := make([]usecases.SkuId, len(items))
	for i, item := range items {
		skuIds[i] = item.SkuId
	}
	return skuIds
}
