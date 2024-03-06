package models

type CartContent struct {
	totalPrice Price
	items      []CartItemInfo
}

func NewCartContent() *CartContent {
	return &CartContent{}
}

func (cc *CartContent) Add(it CartItemInfo) {
	cc.items = append(cc.items, it)
	cc.totalPrice += it.ProductInfo.Price * Price(it.CartItem.Count)
}

func (cc *CartContent) GetPrice() Price {
	return cc.totalPrice
}

func (cc *CartContent) GetItems() []CartItemInfo {
	return cc.items
}
