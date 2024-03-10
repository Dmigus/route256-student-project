package list

type listResponseItem struct {
	SkuId int64  `json:"sku_id"`
	Name  string `json:"name"`
	Count uint16 `json:"count"`
	Price uint32 `json:"price"`
}

type listResponse struct {
	Items      []listResponseItem `json:"items"`
	TotalPrice uint32             `json:"total_price"`
}
