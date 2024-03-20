package loms

import (
	"context"
	"route256.ozon.ru/project/cart/internal/clients/loms/protoc/v1"
)

func (L *Client) GetNumberOfItemInStocks(ctx context.Context, skuId int64) (uint64, error) {
	req := &v1.ListStocksInfoRequest{Sku: uint32(skuId)}
	response, err := L.client.ListStocksInfo(ctx, req)
	if err != nil {
		return 0, err
	}
	return response.Count, nil
}
