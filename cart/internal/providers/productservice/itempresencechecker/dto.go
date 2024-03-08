package itempresencechecker

type listSkusRequest struct {
	Token         string `json:"token"`
	StartAfterSku uint32 `json:"startAfterSku"`
	Count         uint32 `json:"count"`
}

func (l *listSkusRequest) SetToken(token string) {
	l.Token = token
}

type listSkusResponse struct {
	Skus []uint32 `json:"skus"`
}
