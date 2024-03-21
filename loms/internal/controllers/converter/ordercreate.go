package converter

import (
	v1 "route256.ozon.ru/project/loms/internal/controllers/protoc/v1"
	"route256.ozon.ru/project/loms/internal/models"
)

func OrderCreateReqToModel(req *v1.OrderCreateRequest) (int64, []models.OrderItem) {
	items := make([]models.OrderItem, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, ItemToOrderItem(it))
	}
	return req.User, items
}

func IdToOrderCreateResponse(id int64) *v1.OrderId {
	return &v1.OrderId{Id: id}
}
