package cacher

type Option interface {
	apply(c *Cacher)
}

type MaxCacheSizeOption struct {
	size uint
}

func (m MaxCacheSizeOption) apply(c *Cacher) {
	c.cache.maxSize = m.size
}

func WithMaxCacheSize(size uint) MaxCacheSizeOption {
	return MaxCacheSizeOption{size}
}
