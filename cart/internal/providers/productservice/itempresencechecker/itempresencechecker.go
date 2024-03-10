package itempresencechecker

import (
	"context"
	"fmt"
	"math"
	"route256.ozon.ru/project/cart/internal/providers/productservice"
)

var (
	errInvalidSkusArray = fmt.Errorf("no list sku in response")
	errSkuIdIsNotUInt32 = fmt.Errorf("skuId is not in range UInt32")
)

type callPerformer interface {
	Perform(ctx context.Context, method string, reqBody productservice.RequestWithSettableToken, respBody any) error
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
	reqBody := listSkusRequest{
		StartAfterSku: uint32(skuId - 1),
		Count:         1,
	}
	var respDTO listSkusResponse
	err := ipc.rcPerformer.Perform(ctx, "list_skus", &reqBody, &respDTO)
	if err != nil {
		return false, err
	}
	err = validateListSkusResponse(respDTO)
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

func validateListSkusResponse(resp listSkusResponse) error {
	if resp.Skus == nil {
		return errInvalidSkusArray
	}
	return nil
}
