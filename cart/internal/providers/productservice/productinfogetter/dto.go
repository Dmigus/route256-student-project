package productinfogetter

type getProductRequest struct {
	Token string `json:"token"`
	Sku   uint32 `json:"sku"`
}

func (l *getProductRequest) SetToken(token string) {
	l.Token = token
}

type getProductResponse struct {
	Name  *string `json:"name"`
	Price *uint32 `json:"price"`
}
