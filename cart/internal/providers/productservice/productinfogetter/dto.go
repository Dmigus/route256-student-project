package productinfogetter

// GetProductRequest это data transfer object, который представляет собой запрос на получение информации о продукте
type GetProductRequest struct {
	Token string `json:"token"`
	Sku   uint32 `json:"sku"`
}

// SetToken предназначен для установки токена
func (l *GetProductRequest) SetToken(token string) {
	l.Token = token
}

// GetProductResponse это data transfer object, который представляет собой ответ с информацией о продукте
type GetProductResponse struct {
	Name  *string `json:"name"`
	Price *uint32 `json:"price"`
}
