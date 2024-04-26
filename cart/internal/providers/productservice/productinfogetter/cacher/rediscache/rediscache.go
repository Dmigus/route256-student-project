package rediscache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	redisstore "github.com/eko/gocache/store/redis/v4"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"route256.ozon.ru/project/cart/internal/pkg/logging"
	"route256.ozon.ru/project/cart/internal/providers/productservice/productinfogetter/cacher"
	"time"
)

const defaultTTL = 100 * time.Second

type (
	RedisCache struct {
		cache  *cache.Cache[string]
		logger *zap.Logger
		ttl    time.Duration
	}
)

func NewRedisCache(addr string, opts ...Option) *RedisCache {
	redisStore := redisstore.NewRedis(redis.NewClient(&redis.Options{Addr: addr}))
	myCache := &RedisCache{
		cache:  cache.New[string](redisStore),
		logger: logging.EmptyLogger(),
		ttl:    defaultTTL,
	}
	for _, opt := range opts {
		opt.apply(myCache)
	}
	return myCache
}

func (r *RedisCache) Get(ctx context.Context, key cacher.CacheKey) (cacher.CacheValue, bool) {
	val, err := r.cache.Get(ctx, key)
	if err != nil {
		if !isErrorNotFound(err) {
			r.logger.Error("getting from redis failed", logging.AddTraceFieldsFromCtx(ctx, zap.Error(err))...)
		}
		return cacher.CacheValue{}, false
	}
	value, err := r.unmarshalVal(val)
	if err != nil {
		r.logger.Error("could not unmarshal cached value", logging.AddTraceFieldsFromCtx(ctx, zap.Error(err))...)
		return cacher.CacheValue{}, false
	}
	return value, true
}

func isErrorNotFound(err error) bool {
	_, ok := err.(*store.NotFound)
	return ok
}

func (r *RedisCache) Store(ctx context.Context, key cacher.CacheKey, value cacher.CacheValue) {
	bytes, err := r.marshalVal(value)
	if err != nil {
		r.logger.Error("could not marshal value to cache", logging.AddTraceFieldsFromCtx(ctx, zap.Error(err))...)
		return
	}
	err = r.cache.Set(ctx, key, bytes, store.WithExpiration(r.ttl))
	if err != nil {
		r.logger.Error("storing to redis failed", logging.AddTraceFieldsFromCtx(ctx, zap.Error(err))...)
	}
}

func (r *RedisCache) marshalVal(value cacher.CacheValue) (string, error) {
	dtoVal := dto{Response: value.Response}
	if value.Err != nil {
		errStr := value.Err.Error()
		dtoVal.ErrStr = &errStr
	}
	marshalled, err := json.Marshal(dtoVal)
	return string(marshalled), err
}

func (r *RedisCache) unmarshalVal(data string) (cacher.CacheValue, error) {
	dtoVal := dto{}
	err := json.Unmarshal([]byte(data), &dtoVal)
	if err != nil {
		return cacher.CacheValue{}, err
	}
	val := cacher.CacheValue{Response: dtoVal.Response}
	if dtoVal.ErrStr != nil {
		val.Err = errors.New(*dtoVal.ErrStr)
	}
	return val, nil
}
