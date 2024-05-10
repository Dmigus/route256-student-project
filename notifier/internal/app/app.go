// Package app содержит инициализацию и запуск сервиса, согласно конфигурации
package app

import (
	"context"

	"go.uber.org/zap"
	"route256.ozon.ru/project/notifier/internal/providers/eventhandler"
	"route256.ozon.ru/project/notifier/internal/providers/handlingrunner"
	"route256.ozon.ru/project/notifier/internal/service"
)

// App это приложение, которое инициализирует и запускает notifier
type App struct {
	config  Config
	service *service.Notifier
	runner  *handlingrunner.KafkaConsumerGroupRunner
}

// NewApp возращает App, готовый к запуску
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
	runnerConfig := a.createRunnerConfig()
	handler := eventhandler.NewLoggerToWriter(a.config.Logger)
	runner, err := handlingrunner.NewKafkaConsumerGroupRunner(runnerConfig)
	if err != nil {
		return err
	}
	a.runner = runner
	a.service = service.NewNotifier(handler, a.runner)
	return nil
}

func (a *App) createRunnerConfig() handlingrunner.ConsumerGroupRunnerConfig {
	return handlingrunner.ConsumerGroupRunnerConfig{
		Brokers: a.config.Brokers,
		Topic:   a.config.Topic,
		Logger:  a.config.Logger,
	}
}

// Run запускает notifier. Блокирующий.
func (a *App) Run(ctx context.Context) error {
	return a.service.Run(ctx)
}

func (a *App) Update(config Config) {
	a.config = config
	err := a.runner.Update(a.createRunnerConfig())
	if err != nil {
		logger := a.config.Logger
		if logger != nil {
			logger.Error("error updating consumer group runner", zap.Error(err))
		}
	}
}
