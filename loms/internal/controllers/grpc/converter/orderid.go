package converter

import (
	v1 "route256.ozon.ru/project/loms/internal/pkg/api/loms/v1"
)

func OrderIdToId(orderId *v1.OrderId) int64 {
	return orderId.OrderID
}

// IDToOrderID конверирует orderID в grpc OrderID
func IDToOrderID(id int64) *v1.OrderId {
	return &v1.OrderId{OrderID: id}
}
