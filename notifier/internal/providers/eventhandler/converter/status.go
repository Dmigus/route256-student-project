package converter

import (
	"route256.ozon.ru/project/notifier/internal/pkg/api/loms/v1"
)

func TransportStatusToString(st v1.OrderStatus) string {
	switch st {
	case v1.OrderStatus_ORDER_STATUS_NEW:
		return "New"
	case v1.OrderStatus_ORDER_STATUS_AWAITING_PAYMENT:
		return "Awaiting payment"
	case v1.OrderStatus_ORDER_STATUS_PAYED:
		return "Payed"
	case v1.OrderStatus_ORDER_STATUS_FAILED:
		return "Failed"
	case v1.OrderStatus_ORDER_STATUS_CANCELLED:
		return "Cancelled"
	default:
		return "Undefined"
	}
}
