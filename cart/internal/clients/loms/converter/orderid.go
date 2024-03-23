package converter

import v1 "route256.ozon.ru/project/cart/internal/clients/loms/protoc/v1"

func OrderIdToId(orderId *v1.OrderId) int64 {
	return orderId.OrderID
}
