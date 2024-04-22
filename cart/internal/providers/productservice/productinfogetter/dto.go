package productinfogetter

type GetProductRequest struct {
	Token string `json:"token"`
	Sku   uint32 `json:"sku"`
}

func (l *GetProductRequest) SetToken(token string) {
	l.Token = token
}

type GetProductResponse struct {
	Name  *string `json:"name"`
	Price *uint32 `json:"price"`
}
