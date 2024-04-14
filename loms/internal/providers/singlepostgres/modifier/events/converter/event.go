package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"route256.ozon.ru/project/loms/internal/models"
	"route256.ozon.ru/project/loms/internal/pkg/api/loms/v1"
	"time"
)

func NewEventMessage(order *models.Order, dt time.Time) *v1.ChangeOrderStatusEvent {
	return &v1.ChangeOrderStatusEvent{
		OrderID:  &v1.OrderId{OrderID: order.Id()},
		Datetime: timestamppb.New(dt),
		UserID:   order.UserId,
		Status:   ModelOrderStatusToTransportStatus(order.Status),
	}
}
