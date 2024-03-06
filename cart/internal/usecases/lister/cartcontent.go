package lister

import "route256.ozon.ru/project/cart/internal/usecases"

type ItemInfo struct {
	SkuId usecases.SkuId
	Name  string
	Count usecases.ItemCount
	Price usecases.Price
}

type CartContent struct {
	TotalPrice usecases.Price
	Items      []ItemInfo
}

func (cc *CartContent) addItem(it ItemInfo) {
	cc.Items = append(cc.Items, it)
	cc.TotalPrice += it.Price * usecases.Price(it.Count)
}
