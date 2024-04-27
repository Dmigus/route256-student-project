package rediscache

import (
	"time"

	"go.uber.org/zap"
)

// Option это интерфейс для объекта, который может настравивать RedisCache
type Option interface {
	apply(cache *RedisCache)
}

type optionFunc func(cache *RedisCache)

func (o optionFunc) apply(cache *RedisCache) {
	o(cache)
}

// WithLogger устанавливает логер для RedisCache
func WithLogger(logger *zap.Logger) Option {
	return optionFunc(func(cache *RedisCache) {
		cache.logger = logger
	})
}

// WithTTL устанавливает time to live каждой записи, которая сохраняется в кэш
func WithTTL(dur time.Duration) Option {
	return optionFunc(func(cache *RedisCache) {
		cache.ttl = dur
	})
}
