package converter

import (
	v1 "route256.ozon.ru/project/cart/internal/clients/loms/protoc/v1"
	"route256.ozon.ru/project/cart/internal/models"
)

func ModelsToOrderCreateRequest(userId int64, items []models.CartItem) *v1.OrderCreateRequest {
	orderItems := make([]*v1.Item, 0, len(items))
	for _, it := range items {
		orderItems = append(orderItems, CartItemsToOrderItem(it))
	}
	return &v1.OrderCreateRequest{User: userId, Items: orderItems}
}
