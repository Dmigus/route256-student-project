package itempresencechecker

type ListSkusRequest struct {
	Token         string `json:"token"`
	StartAfterSku uint32 `json:"startAfterSku"`
	Count         uint32 `json:"count"`
}

func (l *ListSkusRequest) SetToken(token string) {
	l.Token = token
}

type ListSkusResponse struct {
	Skus []uint32 `json:"skus"`
}
