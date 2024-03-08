package models

type CartContent struct {
	totalPrice uint32
	items      []CartItemInfo
}

func NewCartContent() *CartContent {
	return &CartContent{}
}

func (cc *CartContent) Add(it CartItemInfo) {
	cc.items = append(cc.items, it)
	cc.totalPrice += it.ProductInfo.Price * uint32(it.CartItem.Count)
}

func (cc *CartContent) GetPrice() uint32 {
	return cc.totalPrice
}

func (cc *CartContent) GetItems() []CartItemInfo {
	return cc.items
}
