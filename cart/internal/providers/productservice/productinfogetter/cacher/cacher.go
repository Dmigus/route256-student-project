package cacher

import (
	"context"
	"errors"
	"route256.ozon.ru/project/cart/internal/providers/productservice"
	"route256.ozon.ru/project/cart/internal/providers/productservice/productinfogetter"
)

type (
	CacheKey struct {
		Method  string
		Request productinfogetter.GetProductRequest
	}
	CacheValue struct {
		Response productinfogetter.GetProductResponse
		Err      error
	}
	callPerformer interface {
		Perform(ctx context.Context, method string, reqBody productservice.RequestWithSettableToken) (*productinfogetter.GetProductResponse, error)
	}
	cache interface {
		Get(context.Context, CacheKey) (CacheValue, bool)
		Store(context.Context, CacheKey, CacheValue)
	}
	counter interface {
		Inc()
	}
	Cacher struct {
		rcPerformer                       callPerformer
		cache                             cache
		coordinator                       *execOnceCoordinator
		cacheHitCounter, cacheMissCounter counter
	}
)

func NewCacher(rcPerformer callPerformer, cache cache, cacheHitCounter, cacheMissCounter counter) *Cacher {
	coordinator := newExecOnceCoordinator()
	return &Cacher{rcPerformer: rcPerformer, cache: cache, coordinator: coordinator, cacheHitCounter: cacheHitCounter, cacheMissCounter: cacheMissCounter}
}

func (c *Cacher) Perform(ctx context.Context, method string, reqBody productservice.RequestWithSettableToken) (*productinfogetter.GetProductResponse, error) {
	requestStruct := *reqBody.(*productinfogetter.GetProductRequest)
	key := CacheKey{Method: method, Request: requestStruct}
	result, present := c.cache.Get(ctx, key)
	if !present {
		c.cacheMissCounter.Inc()
		result = c.performExecAndSave(ctx, key)
	} else {
		c.cacheHitCounter.Inc()
	}
	if result.Err != nil {
		return nil, result.Err
	}
	return &result.Response, nil
}

func (c *Cacher) performExecAndSave(ctx context.Context, k CacheKey) CacheValue {
	execOnceAmongGroup := c.coordinator.getExecutor(k)
	funcToExec := c.getPerformAndSaveFunc(ctx, k)
	return execOnceAmongGroup.Execute(funcToExec)
}

func (c *Cacher) getPerformAndSaveFunc(ctx context.Context, k CacheKey) funcToBeExecutedAtMostOnce {
	return func() CacheValue {
		response, err := c.rcPerformer.Perform(ctx, k.Method, &k.Request)
		val := CacheValue{Response: *response, Err: err}
		if !errors.Is(err, context.Canceled) {
			c.cache.Store(ctx, k, val)
		}
		return val
	}
}
