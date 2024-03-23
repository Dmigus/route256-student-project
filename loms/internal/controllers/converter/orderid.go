package converter

import v1 "route256.ozon.ru/project/loms/internal/controllers/protoc/v1"

func OrderIdToId(orderId *v1.OrderId) int64 {
	return orderId.OrderID
}
