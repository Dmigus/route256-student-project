package shardmanager

import "github.com/jackc/pgx/v5/pgxpool"

// Shard представляет шард
type Shard struct {
	master, replica *pgxpool.Pool
}

// NewShard создаёт новый шард с двумя коннектами master replica
func NewShard(master, replica *pgxpool.Pool) *Shard {
	return &Shard{master: master, replica: replica}
}

// Master возвращает соединение к мастер узлу шарда
func (s Shard) Master() *pgxpool.Pool {
	return s.master
}

// Replica возвращает соединение к слейв узлу шарда
func (s Shard) Replica() *pgxpool.Pool {
	return s.replica
}
