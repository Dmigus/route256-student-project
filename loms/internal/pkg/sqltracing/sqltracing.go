// Package sqltracing содержит функции для трейсинга sql запросов
package sqltracing

import (
	"context"
	"fmt"
	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateConnToPostgres создаёт pgxpool.Pool с настроенным трейсингом для каждого запроса
func CreateConnToPostgres(dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}
	cfg.ConnConfig.Tracer = otelpgx.NewTracer(otelpgx.WithTrimSQLInSpanName())
	conn, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}
	if err = conn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}
	return conn, nil
}
