package lister

import (
	"route256.ozon.ru/project/cart/internal/service"
)

type ItemInfo struct {
	SkuId service.SkuId
	Name  string
	Count service.ItemCount
	Price service.Price
}

type CartContent struct {
	totalPrice service.Price
	items      []ItemInfo
}

func (cc *CartContent) addItem(it ItemInfo) {
	cc.items = append(cc.items, it)
	cc.totalPrice += it.Price * service.Price(it.Count)
}
