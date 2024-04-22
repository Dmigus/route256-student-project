package cacher

import (
	"context"
	"errors"
	"route256.ozon.ru/project/cart/internal/providers/productservice"
	"route256.ozon.ru/project/cart/internal/providers/productservice/productinfogetter"
)

type (
	callPerformer interface {
		Perform(ctx context.Context, method string, reqBody productservice.RequestWithSettableToken) (*productinfogetter.GetProductResponse, error)
	}
	Cacher struct {
		rcPerformer callPerformer
		cache       *cache
	}
)

func NewCacher(rcPerformer callPerformer, opts ...Option) *Cacher {
	cacher := &Cacher{rcPerformer: rcPerformer, cache: newCache()}
	for _, opt := range opts {
		opt.apply(cacher)
	}
	return cacher
}

func (c *Cacher) Perform(ctx context.Context, method string, reqBody productservice.RequestWithSettableToken) (*productinfogetter.GetProductResponse, error) {
	requestStruct := *reqBody.(*productinfogetter.GetProductRequest)
	key := key{method: method, request: requestStruct}
	cached, present := c.cache.Get(key)
	if present {
		if cached.err != nil {
			return nil, cached.err
		}
		return &cached.response, nil
	}
	response, err := c.rcPerformer.Perform(ctx, method, reqBody)
	if !errors.Is(err, context.Canceled) {
		c.cache.Insert(key, value{response: *response, err: err})
	}
	return response, err
}
