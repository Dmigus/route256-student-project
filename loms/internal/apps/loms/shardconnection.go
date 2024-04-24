package loms

import (
	"route256.ozon.ru/project/loms/internal/apps"
	"route256.ozon.ru/project/loms/internal/pkg/sqltracing"
	"route256.ozon.ru/project/loms/internal/providers/singlepostgres/shardmanager"
)

type ShardConfig struct {
	Master  apps.PostgresConnectConfig `json:"Master"`
	Replica apps.PostgresConnectConfig `json:"Replica"`
}

func newShardManager(config []ShardConfig) (*shardmanager.Manager, error) {
	connections := make([]shardmanager.ShardConnection, 0, len(config))
	for _, conf := range config {
		conn, err := newShardConnection(conf)
		if err != nil {
			return nil, err
		}
		connections = append(connections, *conn)
	}
	return shardmanager.New(connections)
}

func newShardConnection(config ShardConfig) (*shardmanager.ShardConnection, error) {
	connMaster, err := sqltracing.CreateConnToPostgres(config.Master.GetPostgresDSN())
	if err != nil {
		return nil, err
	}
	connReplica, err := sqltracing.CreateConnToPostgres(config.Replica.GetPostgresDSN())
	if err != nil {
		return nil, err
	}
	return shardmanager.NewShardConnection(connMaster, connReplica), nil
}
