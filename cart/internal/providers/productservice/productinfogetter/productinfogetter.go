package productinfogetter

import (
	"context"
	"errors"
	"fmt"
	"math"
	"route256.ozon.ru/project/cart/internal/models"
	"route256.ozon.ru/project/cart/internal/providers/productservice"
)

var (
	errInvalidPrice       = fmt.Errorf("returned price is not valid")
	errInvalidProductName = fmt.Errorf("returned name is not valid")
	errSkuIdIsNotUInt32   = fmt.Errorf("skuId is not in range UInt32")
)

type callPerformer interface {
	Perform(ctx context.Context, method string, reqBody productservice.RequestWithSettableToken, respBody any) error
}

// ProductInfoGetter прдназначен для возвращения информации о продуктах
type ProductInfoGetter struct {
	rcPerformer callPerformer
}

func NewProductInfoGetter(rcPerformer callPerformer) *ProductInfoGetter {
	return &ProductInfoGetter{
		rcPerformer: rcPerformer,
	}
}

// GetProductsInfo принимает ИД товаров и возвращет их название и цену в том же порядке, как было в skuIds.
func (pig *ProductInfoGetter) GetProductsInfo(ctx context.Context, skuIds []int64) ([]models.ProductInfo, error) {
	prodInfos := make([]models.ProductInfo, 0, len(skuIds))
	for _, skuId := range skuIds {
		prodInfo, err := pig.getProductInfo(ctx, skuId)
		if err != nil {
			return nil, err
		}
		prodInfos = append(prodInfos, prodInfo)
	}
	return prodInfos, nil
}

func (pig *ProductInfoGetter) getProductInfo(ctx context.Context, skuId int64) (models.ProductInfo, error) {
	if err := pig.checkSkuId(skuId); err != nil {
		return models.ProductInfo{}, errSkuIdIsNotUInt32
	}
	reqBody := getProductRequest{
		Sku: uint32(skuId),
	}
	var respDTO getProductResponse
	err := pig.rcPerformer.Perform(ctx, "get_product", &reqBody, &respDTO)
	if err != nil {
		return models.ProductInfo{}, err
	}
	err = validateGetProductResponse(respDTO)
	if err != nil {
		return models.ProductInfo{}, err
	}
	return models.ProductInfo{
		Name:  *respDTO.Name,
		Price: *respDTO.Price,
	}, nil
}

func (pig *ProductInfoGetter) checkSkuId(skuId int64) error {
	if skuId < 0 || skuId > math.MaxUint32 {
		return errSkuIdIsNotUInt32
	}
	return nil
}

func validateGetProductResponse(resp getProductResponse) error {
	var err error
	if resp.Name == nil || !models.IsStringValidName(*resp.Name) {
		err = errors.Join(err, errInvalidProductName)
	}
	if resp.Price == nil {
		err = errors.Join(err, errInvalidPrice)
	}
	return err
}
