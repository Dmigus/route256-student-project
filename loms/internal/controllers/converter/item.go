package converter

import (
	v1 "route256.ozon.ru/project/loms/internal/controllers/protoc/v1"
	"route256.ozon.ru/project/loms/internal/models"
)

func OrderItemToItem(item models.OrderItem) *v1.Item {
	return &v1.Item{Sku: uint32(item.SkuId), Count: uint32(item.Count)}
}

func ItemToOrderItem(it *v1.Item) models.OrderItem {
	return models.OrderItem{SkuId: int64(it.Sku), Count: uint16(it.Count)}
}
