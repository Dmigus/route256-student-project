package lister

import (
	"route256.ozon.ru/project/cart/internal/service/modifier"
)

type ItemInfo struct {
	SkuId modifier.SkuId
	Name  string
	Count uint16
	Price uint32
}

type CartContent struct {
	totalPrice uint32
	items      []ItemInfo
}
