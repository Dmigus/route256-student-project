package rediscache

import (
	"go.uber.org/zap"
	"time"
)

type Option interface {
	apply(cache *RedisCache)
}

type optionFunc func(cache *RedisCache)

func (o optionFunc) apply(cache *RedisCache) {
	o(cache)
}

func WithLogger(logger *zap.Logger) Option {
	return optionFunc(func(cache *RedisCache) {
		cache.logger = logger
	})
}

func WithTTL(dur time.Duration) Option {
	return optionFunc(func(cache *RedisCache) {
		cache.ttl = dur
	})
}
