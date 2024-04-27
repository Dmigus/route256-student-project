package shardmanager

// Option это интерфейс для объекта, который может настравивать Manager
type Option interface {
	apply(*Manager)
}

type optionFunc func(*Manager)

func (o optionFunc) apply(c *Manager) {
	o(c)
}

// WithHashFn задаёт хэш функцию дял использования ShardManager
func WithHashFn(hf HashFn) Option {
	return optionFunc(func(manager *Manager) {
		manager.fn = hf
	})
}

// WithDefaultShard задаёт шард, который будет считаться "дефолтным"
func WithDefaultShard(sc Shard) Option {
	return optionFunc(func(manager *Manager) {
		manager.defaultShard = sc
	})
}
