// Package outboxsender содержит приложение, которое выполняет выгрузку из outbox в брокер
package outboxsender

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"route256.ozon.ru/project/loms/internal/providers/kafka"
	"route256.ozon.ru/project/loms/internal/providers/singlepostgres"
	"route256.ozon.ru/project/loms/internal/providers/singlepostgres/modifier"
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
	pusher, err := kafka.NewProducer(a.config.Kafka.Brokers, a.config.Kafka.Topic)
	if err != nil {
		return err
	}
	connOutbox := createConnToPostgres(a.config.Outbox.GetPostgresDSN())
	txM := singlepostgres.NewTxManagerOne(connOutbox, func(conn pgx.Tx) outboxsender.Outbox {
		return modifier.NewEvents(conn)
	})
	a.service = outboxsender.NewService(txM, pusher, time.Duration(a.config.BatchInterval)*time.Second, a.config.BatchSize)
	return nil
}

// Run это блокирующий запуск приложения App с временем жизни, определяемым ctx
func (a *App) Run(ctx context.Context) {
	a.service.Run(ctx)
}

func createConnToPostgres(dsn string) *pgxpool.Pool {
	conn, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
