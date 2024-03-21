package converter

import v1 "route256.ozon.ru/project/loms/internal/controllers/protoc/v1"

func ListStocksInfoRequestToSkuId(req *v1.StocksInfoRequest) int64 {
	return int64(req.Sku)
}

func CountToStocksInfoResponse(count uint64) *v1.StocksInfoResponse {
	return &v1.StocksInfoResponse{Count: count}
}
