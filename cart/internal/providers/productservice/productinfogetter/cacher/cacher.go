package cacher

import (
	"context"
	"errors"
	"route256.ozon.ru/project/cart/internal/providers/productservice"
	"route256.ozon.ru/project/cart/internal/providers/productservice/productinfogetter"
)

type (
	CacheKey struct {
		method  string
		request productinfogetter.GetProductRequest
	}
	CacheValue struct {
		response productinfogetter.GetProductResponse
		err      error
	}
	callPerformer interface {
		Perform(ctx context.Context, method string, reqBody productservice.RequestWithSettableToken) (*productinfogetter.GetProductResponse, error)
	}
	cache interface {
		Get(CacheKey) (CacheValue, bool)
		Store(CacheKey, CacheValue)
	}
	Cacher struct {
		rcPerformer callPerformer
		cache       cache
		coordinator *execOnceCoordinator
	}
)

func NewCacher(rcPerformer callPerformer, cache cache) *Cacher {
	coordinator := newExecOnceCoordinator()
	return &Cacher{rcPerformer: rcPerformer, cache: cache, coordinator: coordinator}
}

func (c *Cacher) Perform(ctx context.Context, method string, reqBody productservice.RequestWithSettableToken) (*productinfogetter.GetProductResponse, error) {
	requestStruct := *reqBody.(*productinfogetter.GetProductRequest)
	key := CacheKey{method: method, request: requestStruct}
	result, present := c.cache.Get(key)
	if !present {
		result = c.performExecAndSave(ctx, key)
	}
	if result.err != nil {
		return nil, result.err
	}
	return &result.response, nil
}

func (c *Cacher) performExecAndSave(ctx context.Context, k CacheKey) CacheValue {
	execOnceAmongGroup := c.coordinator.getExecutor(k)
	funcToExec := c.getPerformAndSaveFunc(ctx, k)
	return execOnceAmongGroup.Execute(funcToExec)
}

func (c *Cacher) getPerformAndSaveFunc(ctx context.Context, k CacheKey) funcToBeExecutedAtMostOnce {
	return func() CacheValue {
		response, err := c.rcPerformer.Perform(ctx, k.method, &k.request)
		val := CacheValue{response: *response, err: err}
		if !errors.Is(err, context.Canceled) {
			c.cache.Store(k, val)
		}
		return val
	}
}
