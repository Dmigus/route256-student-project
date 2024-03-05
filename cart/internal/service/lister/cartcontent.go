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
	TotalPrice service.Price
	Items      []ItemInfo
}

func (cc *CartContent) addItem(it ItemInfo) {
	cc.Items = append(cc.Items, it)
	cc.TotalPrice += it.Price * service.Price(it.Count)
}
