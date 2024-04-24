package shardmanager

type Option interface {
	apply(*Manager)
}

type optionFunc func(*Manager)

func (o optionFunc) apply(c *Manager) {
	o(c)
}

func WithHashFn(hf HashFn) Option {
	return optionFunc(func(manager *Manager) {
		manager.fn = hf
	})
}

func WithDefaultShard(sc ShardConnection) Option {
	return optionFunc(func(manager *Manager) {
		manager.defaultShard = sc
	})
}
