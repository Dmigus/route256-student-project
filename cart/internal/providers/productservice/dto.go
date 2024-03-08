package productservice

type listSkusRequest struct {
	Token         string `json:"token"`
	StartAfterSku uint32 `json:"startAfterSku"`
	Count         uint32 `json:"count"`
}

type listSkusResponse struct {
	Skus *[]uint32 `json:"skus"`
}

type getProductRequest struct {
	Token string `json:"token"`
	Sku   uint32 `json:"sku"`
}

type getProductResponse struct {
	Name  *string `json:"name"`
	Price *uint32 `json:"price"`
}
