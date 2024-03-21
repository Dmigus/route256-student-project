package converter

import v1 "route256.ozon.ru/project/cart/internal/clients/loms/protoc/v1"

func SkuIdToStocksInfoRequest(skuId int64) *v1.StocksInfoRequest {
	return &v1.StocksInfoRequest{Sku: uint32(skuId)}
}

func StocksInfoResponseToCount(resp *v1.StocksInfoResponse) uint64 {
	return resp.Count
}
