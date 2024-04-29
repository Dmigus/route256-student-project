// Package outboxsender содержит приложение, которое выполняет выгрузку из outbox в брокер
package outboxsender

import (
	"context"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"route256.ozon.ru/project/loms/internal/pkg/sqlmetrics"
	"route256.ozon.ru/project/loms/internal/pkg/sqltracing"
	"route256.ozon.ru/project/loms/internal/providers/kafka"
	"route256.ozon.ru/project/loms/internal/providers/multipostgres"
	"route256.ozon.ru/project/loms/internal/providers/multipostgres/modifiers/events"
	"route256.ozon.ru/project/loms/internal/services/outboxsender"
)

var bucketsForRequestDuration = []float64{0.001, 0.005, 0.01, 0.02, 0.05, 0.1, 0.2, 0.5, 1}

// App это приложение для отправителя из outbox в топик кафки
type App struct {
	config   Config
	services []*outboxsender.Service
}

// NewApp создаёт и готовит к запуску новое приложение App, согласно переданной конфигурации. Возвращает ошибку, если не удалось сконфигурировать.
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
	sqlDurationRecorder, err := a.initSQLMetric()
	if err != nil {
		return err
	}
	for _, outboxShard := range a.config.Outbox {
		pusher, err := kafka.NewSender(a.config.Kafka.Brokers, a.config.Kafka.Topic)
		if err != nil {
			return err
		}
		connOutbox, err := sqltracing.CreateConnToPostgres(outboxShard.GetPostgresDSN())
		if err != nil {
			return err
		}
		txManagerShard := multipostgres.NewTxManager1(func(creator multipostgres.TransactionCreator) outboxsender.Outbox {
			return events.NewEventsToPull(creator, connOutbox, sqlDurationRecorder)
		})
		serviceShard := outboxsender.NewService(txManagerShard, pusher, time.Duration(a.config.BatchInterval)*time.Second, a.config.BatchSize, a.config.Logger)
		a.services = append(a.services, serviceShard)
	}
	return nil
}

func (a *App) initSQLMetric() (*sqlmetrics.SQLRequestDuration, error) {
	responseTime := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "outboxsender",
		Name:      "sql_request_duration_seconds",
		Help:      "Response time distribution made to PostgreSQL",
		Buckets:   bucketsForRequestDuration,
	},
		[]string{sqlmetrics.TableLabel, sqlmetrics.CategoryLabel, sqlmetrics.ErrLabel},
	)
	err := a.config.MetricsRegisterer.Register(responseTime)
	if err != nil {
		return nil, err
	}
	return sqlmetrics.NewSQLRequestDuration(responseTime), nil
}

// Run это блокирующий запуск приложения App с временем жизни, определяемым ctx
func (a *App) Run(ctx context.Context) {
	wg := sync.WaitGroup{}
	for _, service := range a.services {
		wg.Add(1)
		go func() {
			defer wg.Done()
			service.Run(ctx)
		}()
	}
	wg.Wait()
}
