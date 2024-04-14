package converter

import (
	"route256.ozon.ru/project/loms/internal/models"
	"route256.ozon.ru/project/loms/internal/pkg/api/loms/v1"
)

func OrderCreateReqToModel(req *v1.OrderCreateRequest) (int64, []models.OrderItem) {
	items := make([]models.OrderItem, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, ItemToOrderItem(it))
	}
	return req.User, items
}

func IdToOrderCreateResponse(id int64) *v1.OrderId {
	return &v1.OrderId{OrderID: id}
}
