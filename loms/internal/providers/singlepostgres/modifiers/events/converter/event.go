// Package converter содержит функции для работы со структурами, сгенерированными из proto файла, в рамках пакета events
package converter

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"route256.ozon.ru/project/loms/internal/models"
	v1 "route256.ozon.ru/project/loms/internal/pkg/api/loms/v1"
)

// NewEventMessage создаёт экземпляр ChangeOrderStatusEvent из заказа order и текущего времени dt
func NewEventMessage(order *models.Order, dt time.Time) *v1.ChangeOrderStatusEvent {
	return &v1.ChangeOrderStatusEvent{
		OrderID:  &v1.OrderId{OrderID: order.Id()},
		Datetime: timestamppb.New(dt),
		UserID:   order.UserId,
		Status:   ModelOrderStatusToTransportStatus(order.Status),
	}
}
