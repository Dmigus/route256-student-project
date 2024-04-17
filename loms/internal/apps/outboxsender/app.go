// Package outboxsender содержит приложение, которое выполняет выгрузку из outbox в брокер
package outboxsender

import (
	"context"
	"time"

	"route256.ozon.ru/project/loms/internal/providers/singlepostgres/modifiers/events"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"route256.ozon.ru/project/loms/internal/providers/kafka"
	"route256.ozon.ru/project/loms/internal/providers/singlepostgres"
	"route256.ozon.ru/project/loms/internal/services/outboxsender"
)

// App это приложение для отправителя из outbox в топик кафки
type App struct {
	config  Config
	service *outboxsender.Service
}

// NewApp создаёт и готовит к азпуску новое приложение App, согласно переданной конфигурации. Возвращает ошибку, если не удалось сконфигурировать.
func NewApp(config Config) (*App, error) {
	app := &App{
		config: config,
	}
	if err := app.init(); err != nil {
		return nil, err
	}
	return app, nil
}

func (a *App) init() error {
	pusher, err := kafka.NewSender(a.config.Kafka.Brokers, a.config.Kafka.Topic)
	if err != nil {
		return err
	}
	connOutbox, err := createConnToPostgres(a.config.Outbox.GetPostgresDSN())
	if err != nil {
		return err
	}
	txM := singlepostgres.NewTxManagerOne(connOutbox, func(conn pgx.Tx) outboxsender.Outbox {
		return events.NewEvents(conn)
	})
	a.service = outboxsender.NewService(txM, pusher, time.Duration(a.config.BatchInterval)*time.Second, a.config.BatchSize)
	return nil
}

// Run это блокирующий запуск приложения App с временем жизни, определяемым ctx
func (a *App) Run(ctx context.Context) {
	a.service.Run(ctx)
}

func createConnToPostgres(dsn string) (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return conn, nil
}
