// Package shardmanager содержит реализацию шард менеджера
package shardmanager

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/spaolacci/murmur3"
)

// BucketsNum это количество бакетов, в которыми работает весь сервис
const BucketsNum = 1000

type (
	// VShard это виртуальный шард, которому принадлежит тот или иной заказ
	VShard int64
	// ShardHash хэш бакета
	ShardHash uint32
	// HashFn хеш функция, используемая для преобразования бакета в число
	HashFn func(VShard) ShardHash
	// Manager это структура, ответсвенная за хранение и выдачу шардов
	Manager struct {
		fn           HashFn
		shards       []Shard
		defaultShard Shard
	}
)

var errShardsEmpty = errors.New("shards must be not empty")

// New создаёт новый Manager
func New(shards []Shard, opts ...Option) (*Manager, error) {
	if len(shards) == 0 {
		return nil, errShardsEmpty
	}
	m := &Manager{
		fn:           Murmur3HashFn(),
		shards:       shards,
		defaultShard: shards[0],
	}
	for _, opt := range opts {
		opt.apply(m)
	}
	return m, nil
}

// GetShard возвращает шард, относящийся к виртуальному
func (m *Manager) GetShard(key VShard) Shard {
	shardHash := m.fn(key)
	return m.mapHashtoShard(shardHash)
}

func (m *Manager) mapHashtoShard(hash ShardHash) Shard {
	return m.shards[int(hash)%len(m.shards)]
}

// GetDefaultShard возвращает "дефолтный" шард
func (m *Manager) GetDefaultShard() Shard {
	return m.defaultShard
}

// ShardNum возвращает количество шардов всего
func (m *Manager) ShardNum() int {
	return len(m.shards)
}

// GetShardByInd возвращает шард по порядковому номеру ind
func (m *Manager) GetShardByInd(ind int) Shard {
	return m.shards[ind]
}

// Murmur3HashFn возвращает хеш функцию, которая вычисляет по номеру бакета вычисляет хэш
func Murmur3HashFn() HashFn {
	// всегда инициируется with seed = 0
	hasher := murmur3.New32()
	return func(key VShard) ShardHash {
		defer hasher.Reset()
		buf := new(bytes.Buffer)
		_ = binary.Write(buf, binary.LittleEndian, key)
		_, _ = hasher.Write(buf.Bytes())
		return ShardHash(hasher.Sum32())
	}
}
