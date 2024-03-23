package converter

import (
	"route256.ozon.ru/project/loms/internal/controllers/grpc/protoc/v1"
)

func OrderIdToId(orderId *v1.OrderId) int64 {
	return orderId.OrderID
}
