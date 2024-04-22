package inmemorycache

type Option interface {
	apply(c *InMemoryCache)
}

type MaxCacheSizeOption struct {
	size uint
}

func (m MaxCacheSizeOption) apply(c *InMemoryCache) {
	c.maxSize = m.size
}

func WithMaxCacheSize(size uint) MaxCacheSizeOption {
	return MaxCacheSizeOption{size}
}
