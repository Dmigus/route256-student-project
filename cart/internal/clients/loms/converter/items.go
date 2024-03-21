package converter

import (
	v1 "route256.ozon.ru/project/cart/internal/clients/loms/protoc/v1"
	"route256.ozon.ru/project/cart/internal/models"
)

func CartItemsToOrderItem(it models.CartItem) *v1.Item {
	return &v1.Item{
		Sku:   uint32(it.SkuId),
		Count: uint32(it.Count),
	}
}
