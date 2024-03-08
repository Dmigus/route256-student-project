package models

type CartItem struct {
	SkuId int64
	Count uint16
}

func IsNumberValidCount(num uint16) bool {
	return num > 0
}
