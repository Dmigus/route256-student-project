package converter

import (
	"route256.ozon.ru/project/cart/internal/models"
	v1 "route256.ozon.ru/project/cart/internal/providers/loms/client/protoc/v1"
)

func CartItemsToOrderItem(it models.CartItem) *v1.Item {
	return &v1.Item{
		Sku:   uint32(it.SkuId),
		Count: uint32(it.Count),
	}
}
