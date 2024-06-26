package itempresencechecker

import (
	"context"
	"errors"
	"math"
	"route256.ozon.ru/project/cart/internal/providers/productservice"
)

var (
	errInvalidSkusArray = errors.New("no list sku in response")
	errSkuIdIsNotUInt32 = errors.New("skuId is not in range UInt32")
)

type callPerformer interface {
	Perform(ctx context.Context, method string, reqBody productservice.RequestWithSettableToken) (*ListSkusResponse, error)
}

// ItemPresenceChecker предназначен для проверки существования товара
type ItemPresenceChecker struct {
	rcPerformer callPerformer
}

func NewItemPresenceChecker(rcPerformer callPerformer) *ItemPresenceChecker {
	return &ItemPresenceChecker{
		rcPerformer: rcPerformer,
	}
}

// IsItemPresent принимает ИД товара и возращает true, если он существует в "специальном сервисе"
func (ipc *ItemPresenceChecker) IsItemPresent(ctx context.Context, skuId int64) (bool, error) {
	if err := ipc.checkSkuId(skuId); err != nil {
		return false, errSkuIdIsNotUInt32
	}
	reqBody := ListSkusRequest{
		StartAfterSku: uint32(skuId - 1),
		Count:         1,
	}
	respDTO, err := ipc.rcPerformer.Perform(ctx, "list_skus", &reqBody)
	if err != nil {
		return false, err
	}
	err = validateListSkusResponse(*respDTO)
	if err != nil {
		return false, err
	}
	if len(respDTO.Skus) > 0 && respDTO.Skus[0] == uint32(skuId) {
		return true, nil
	}
	return false, nil
}

func (ipc *ItemPresenceChecker) checkSkuId(skuId int64) error {
	if skuId < 0 || skuId > math.MaxUint32 {
		return errSkuIdIsNotUInt32
	}
	return nil
}

func validateListSkusResponse(resp ListSkusResponse) error {
	if resp.Skus == nil {
		return errInvalidSkusArray
	}
	return nil
}
