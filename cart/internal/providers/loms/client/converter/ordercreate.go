package converter

import (
	"route256.ozon.ru/project/cart/internal/models"
	v1 "route256.ozon.ru/project/cart/internal/providers/loms/client/protoc/v1"
)

func ModelsToOrderCreateRequest(userId int64, items []models.CartItem) *v1.OrderCreateRequest {
	orderItems := make([]*v1.Item, 0, len(items))
	for _, it := range items {
		orderItems = append(orderItems, CartItemsToOrderItem(it))
	}
	return &v1.OrderCreateRequest{User: userId, Items: orderItems}
}
