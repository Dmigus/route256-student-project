package productservice

type listSkusRequest struct {
	Token         string `json:"token"`
	StartAfterSku int64  `json:"startAfterSku"`
	Count         int64  `json:"count"`
}

type listSkusResponse struct {
	Skus []int64 `json:"skus"`
}

type getProductRequest struct {
	Token string `json:"token"`
	Sku   int64  `json:"sku"`
}

type getProductResponse struct {
	Name  string `json:"name"`
	Price int64  `json:"price"`
}
