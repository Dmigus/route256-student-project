package converter

import (
	"route256.ozon.ru/project/loms/internal/pkg/api/loms/v1"
)

func OrderIdToId(orderId *v1.OrderId) int64 {
	return orderId.OrderID
}
