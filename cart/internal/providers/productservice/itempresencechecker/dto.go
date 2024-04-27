package itempresencechecker

// ListSkusRequest это data transfer object, который представляет собой запрос на получение информации о продуктах
type ListSkusRequest struct {
	Token         string `json:"token"`
	StartAfterSku uint32 `json:"startAfterSku"`
	Count         uint32 `json:"count"`
}

// SetToken предназначен для установки токена
func (l *ListSkusRequest) SetToken(token string) {
	l.Token = token
}

// ListSkusResponse это data transfer object, который представляет собой ответ с информацией о продуктах
type ListSkusResponse struct {
	Skus []uint32 `json:"skus"`
}
