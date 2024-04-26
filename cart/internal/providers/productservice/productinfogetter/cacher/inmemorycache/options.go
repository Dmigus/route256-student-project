package inmemorycache

type Option interface {
	apply(c *InMemoryCache)
}

type optionFunc func(cache *InMemoryCache)

func (o optionFunc) apply(cache *InMemoryCache) {
	o(cache)
}

func WithMaxCacheSize(size uint) Option {
	return optionFunc(func(cache *InMemoryCache) {
		cache.maxSize = size
	})
}
