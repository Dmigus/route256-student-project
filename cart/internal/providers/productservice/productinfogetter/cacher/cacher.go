package cacher

import (
	"context"
	"errors"
	"route256.ozon.ru/project/cart/internal/providers/productservice"
	"route256.ozon.ru/project/cart/internal/providers/productservice/productinfogetter"
	"time"
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
	summary interface {
		Observe(float64)
	}
	Cacher struct {
		rcPerformer                       callPerformer
		cache                             cache
		coordinator                       *execOnceCoordinator
		cacheHitSummary, cacheMissSummary summary
	}
)

func NewCacher(rcPerformer callPerformer, cache cache, cacheHitSummary, cacheMissSummary summary) *Cacher {
	coordinator := newExecOnceCoordinator()
	return &Cacher{rcPerformer: rcPerformer, cache: cache, coordinator: coordinator, cacheHitSummary: cacheHitSummary, cacheMissSummary: cacheMissSummary}
}

func (c *Cacher) Perform(ctx context.Context, method string, reqBody productservice.RequestWithSettableToken) (*productinfogetter.GetProductResponse, error) {
	requestStruct := *reqBody.(*productinfogetter.GetProductRequest)
	key := CacheKey{Method: method, Request: requestStruct}
	var result CacheValue
	var present bool
	dur := compDuration(func() {
		result, present = c.cache.Get(ctx, key)
	})
	if !present {
		result = c.performExecAndSave(ctx, key)
	} else {
		c.cacheHitSummary.Observe(dur)
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
		var response *productinfogetter.GetProductResponse
		var err error
		dur := compDuration(func() {
			response, err = c.rcPerformer.Perform(ctx, k.Method, &k.Request)
		})
		c.cacheMissSummary.Observe(dur)
		val := CacheValue{Err: err}
		if response != nil {
			val.Response = *response
		}
		if !errors.Is(err, context.Canceled) {
			c.cache.Store(ctx, k, val)
		}
		return val
	}
}

func compDuration(f func()) float64 {
	before := time.Now()
	f()
	after := time.Now()
	return after.Sub(before).Seconds()
}
