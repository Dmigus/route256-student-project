package productservice

import (
	"context"
	"route256.ozon.ru/project/cart/internal/service"
	"route256.ozon.ru/project/cart/internal/service/lister"
)

type ProductService struct {
}

func (p *ProductService) IsItemPresent(ctx context.Context, skuId service.SkuId) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (p *ProductService) GetProductsInfo(ctx context.Context, skuIds []service.SkuId) ([]lister.ProductInfo, error) {
	//TODO implement me
	panic("implement me")
}
