package converter

import (
	v1 "route256.ozon.ru/project/loms/internal/controllers/protoc/v1"
	"route256.ozon.ru/project/loms/internal/models"
)

func OrderToOrderInfoResponse(order *models.Order) *v1.OrderInfoResponse {
	response := &v1.OrderInfoResponse{Items: make([]*v1.Item, 0, len(order.Items))}
	response.Status = ModelOrderStatusToTransportStatus(order.Status)
	response.User = order.UserId
	for _, it := range order.Items {
		response.Items = append(response.Items, OrderItemToItem(it))
	}
	return response
}
