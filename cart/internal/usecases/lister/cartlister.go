package lister

import (
	"context"
	"route256.ozon.ru/project/cart/internal/models"
)

type repository interface {
	GetCart(ctx context.Context, user int64) (*models.Cart, error)
}

type productService interface {
	GetProductsInfo(ctx context.Context, skuIds []int64) ([]models.ProductInfo, error)
}

type CartListerService struct {
	repo           repository
	productService productService
}

func New(repo repository, productService productService) *CartListerService {
	return &CartListerService{repo: repo, productService: productService}
}

func (cl *CartListerService) ListCartContent(ctx context.Context, user int64) (*models.CartContent, error) {
	cart, err := cl.repo.GetCart(ctx, user)
	if err != nil {
		return nil, err
	}
	items := cart.ListItems(ctx)
	skuIds := extractSkuIds(items)
	productInfos, err := cl.productService.GetProductsInfo(ctx, skuIds)
	if err != nil {
		return nil, err
	}
	return createCartContent(items, productInfos), nil
}

func createCartContent(items []models.CartItem, prodInfos []models.ProductInfo) *models.CartContent {
	content := models.NewCartContent()
	for i := range items {
		itInfo := models.CartItemInfo{
			CartItem:    items[i],
			ProductInfo: prodInfos[i],
		}
		content.Add(itInfo)
	}
	return content
}

func extractSkuIds(items []models.CartItem) []int64 {
	skuIds := make([]int64, len(items))
	for i, item := range items {
		skuIds[i] = item.SkuId
	}
	return skuIds
}
