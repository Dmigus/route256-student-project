package loms

import (
	"context"
	"route256.ozon.ru/project/cart/internal/clients/loms/protoc/v1"
	"route256.ozon.ru/project/cart/internal/models"
)

func (L *Client) CreateOrder(ctx context.Context, userId int64, items []models.CartItem) (int64, error) {
	orderItems := make([]*v1.Item, 0, len(items))
	for _, it := range items {
		orderItems = append(orderItems, cartItemsToOrderItem(it))
	}
	request := &v1.CreateOrderRequest{User: userId, Items: orderItems}
	response, err := L.client.CreateOrder(ctx, request)
	if err != nil {
		return 0, err
	}
	return response.Id, nil
}

func cartItemsToOrderItem(it models.CartItem) *v1.Item {
	return &v1.Item{
		Sku:   uint32(it.SkuId),
		Count: uint32(it.Count),
	}
}
