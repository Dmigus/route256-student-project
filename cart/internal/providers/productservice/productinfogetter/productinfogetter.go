package productinfogetter

import (
	"context"
	"errors"
	"math"
	"route256.ozon.ru/project/cart/internal/models"
	"route256.ozon.ru/project/cart/internal/pkg/errorgroup"
	"route256.ozon.ru/project/cart/internal/providers/productservice"
)

var (
	errInvalidPrice       = errors.New("returned price is not valid")
	errInvalidProductName = errors.New("returned name is not valid")
	errSkuIdIsNotUInt32   = errors.New("skuId is not in range UInt32")
)

const remoteMethodName = "get_product"

type callPerformer interface {
	Perform(ctx context.Context, method string, reqBody productservice.RequestWithSettableToken) (*GetProductResponse, error)
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
	prodInfos := make([]models.ProductInfo, len(skuIds))
	errGr, groupCtx := errorgroup.NewErrorGroup(ctx)
	for ind, skuID := range skuIds {
		errGr.Go(func() error {
			prodInfo, err := pig.getProductInfo(groupCtx, skuID)
			if err != nil {
				return err
			}
			prodInfos[ind] = prodInfo
			return nil
		})
	}
	if err := errGr.Wait(); err != nil {
		return nil, err
	}
	return prodInfos, nil
}

func (pig *ProductInfoGetter) getProductInfo(ctx context.Context, skuId int64) (models.ProductInfo, error) {
	if err := pig.checkSkuId(skuId); err != nil {
		return models.ProductInfo{}, errSkuIdIsNotUInt32
	}
	reqBody := GetProductRequest{
		Sku: uint32(skuId),
	}
	respDTO, err := pig.rcPerformer.Perform(ctx, remoteMethodName, &reqBody)
	if err != nil {
		return models.ProductInfo{}, err
	}
	err = validateGetProductResponse(*respDTO)
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

func validateGetProductResponse(resp GetProductResponse) error {
	var err error
	if resp.Name == nil || !models.IsStringValidName(*resp.Name) {
		err = errors.Join(err, errInvalidProductName)
	}
	if resp.Price == nil {
		err = errors.Join(err, errInvalidPrice)
	}
	return err
}
