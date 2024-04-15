package converter

import (
	"route256.ozon.ru/project/loms/internal/models"
	v1 "route256.ozon.ru/project/loms/internal/pkg/api/loms/v1"
)

// ModelOrderStatusToTransportStatus сопоставляет доменный OrderStatus сгенерированному
func ModelOrderStatusToTransportStatus(st models.OrderStatus) v1.OrderStatus {
	switch st {
	case models.New:
		return v1.OrderStatus_ORDER_STATUS_NEW
	case models.AwaitingPayment:
		return v1.OrderStatus_ORDER_STATUS_AWAITING_PAYMENT
	case models.Payed:
		return v1.OrderStatus_ORDER_STATUS_PAYED
	case models.Failed:
		return v1.OrderStatus_ORDER_STATUS_FAILED
	case models.Cancelled:
		return v1.OrderStatus_ORDER_STATUS_CANCELLED
	default:
		return v1.OrderStatus_ORDER_STATUS_UNDEFINED
	}
}
