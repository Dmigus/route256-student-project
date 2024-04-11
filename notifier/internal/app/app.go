package app

import (
	"context"
	"os"
	"route256.ozon.ru/project/notifier/internal/providers/eventhandler"
	"route256.ozon.ru/project/notifier/internal/providers/handlingrunner"
	"route256.ozon.ru/project/notifier/internal/service"
)

type App struct {
	config  Config
	service *service.Notifier
}

// NewApp возращает App, готовый к запуску, либо ошибку
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
	runner, err := handlingrunner.NewKafkaConsumerGroupRunner(a.config.Brokers, a.config.Topic)
	if err != nil {
		return err
	}
	handler := eventhandler.NewLoggerToWriter(os.Stdout)
	a.service = service.NewNotifier(handler, runner)
	return nil
}

func (a *App) Run(ctx context.Context) error {
	return a.service.Run(ctx)
}
