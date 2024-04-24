package shardmanager

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/spaolacci/murmur3"
)

type (
	ShardKey  int64
	ShardHash uint32
	HashFn    func(ShardKey) ShardHash
	Manager   struct {
		fn           HashFn
		shards       []ShardConnection
		defaultShard ShardConnection
	}
)

var errShardsEmpty = errors.New("shards must be not empty")

func New(shards []ShardConnection, opts ...Option) (*Manager, error) {
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

func (m *Manager) GetShard(key ShardKey) ShardConnection {
	shardHash := m.fn(key)
	return m.mapHashtoShard(shardHash)
}

func (m *Manager) mapHashtoShard(hash ShardHash) ShardConnection {
	return m.shards[int(hash)%len(m.shards)]
}

func (m *Manager) GetDefaultShard() ShardConnection {
	return m.defaultShard
}

func Murmur3HashFn() HashFn {
	// всегда инициируется with seed = 0
	hasher := murmur3.New32()
	return func(key ShardKey) ShardHash {
		defer hasher.Reset()
		buf := new(bytes.Buffer)
		_ = binary.Write(buf, binary.LittleEndian, key)
		_, _ = hasher.Write(buf.Bytes())
		return ShardHash(hasher.Sum32())
	}
}
