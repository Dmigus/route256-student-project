package shardmanager

import "github.com/jackc/pgx/v5/pgxpool"

type ShardConnection struct {
	master, replica *pgxpool.Pool
}

func NewShardConnection(master, replica *pgxpool.Pool) *ShardConnection {
	return &ShardConnection{master: master, replica: replica}
}

func (s ShardConnection) Master() *pgxpool.Pool {
	return s.master
}

func (s ShardConnection) Replica() *pgxpool.Pool {
	return s.replica
}
