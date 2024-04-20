// Package outboxsender содержит приложение, которое выполняет выгрузку из outbox в брокер
package outboxsender

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"route256.ozon.ru/project/loms/internal/pkg/sqlmetrics"
	"time"

	"route256.ozon.ru/project/loms/internal/providers/singlepostgres/modifiers/events"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"route256.ozon.ru/project/loms/internal/providers/kafka"
	"route256.ozon.ru/project/loms/internal/providers/singlepostgres"
	"route256.ozon.ru/project/loms/internal/services/outboxsender"
)

var bucketsForRequestDuration = []float64{0.001, 0.005, 0.01, 0.02, 0.05, 0.1, 0.2, 0.5, 1}

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

	responseTime := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "outboxsender",
		Name:      "sql_request_duration_seconds",
		Help:      "Response time distribution made to PostgreSQL",
		Buckets:   bucketsForRequestDuration,
	},
		[]string{sqlmetrics.TableLabel, sqlmetrics.CategoryLabel, sqlmetrics.ErrLabel},
	)
	err = a.config.MetricsRegisterer.Register(responseTime)
	if err != nil {
		return err
	}
	sqlDurationRecorder := sqlmetrics.NewSQLRequestDuration(responseTime)

	txM := singlepostgres.NewTxManagerOne(connOutbox, func(conn pgx.Tx) outboxsender.Outbox {
		return events.NewEvents(conn, sqlDurationRecorder)
	})
	a.service = outboxsender.NewService(txM, pusher, time.Duration(a.config.BatchInterval)*time.Second, a.config.BatchSize, a.config.Logger)
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
