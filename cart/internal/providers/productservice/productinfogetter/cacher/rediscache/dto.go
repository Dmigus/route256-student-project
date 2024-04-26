package rediscache

import "route256.ozon.ru/project/cart/internal/providers/productservice/productinfogetter"

type dto struct {
	Response productinfogetter.GetProductResponse
	ErrStr   *string `json:"ErrStr,omitempty"`
}
