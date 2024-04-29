package inmemorycache

// Option это интерфейс для объекта, который может настравивать inmemorycache
type Option interface {
	apply(c *InMemoryCache)
}

type optionFunc func(cache *InMemoryCache)

func (o optionFunc) apply(cache *InMemoryCache) {
	o(cache)
}

// WithMaxCacheSize устанавливает лимит размера кэша. Если поместить новое значение сверх лимита, то какое-то
func WithMaxCacheSize(size uint) Option {
	return optionFunc(func(cache *InMemoryCache) {
		cache.maxSize = size
	})
}
