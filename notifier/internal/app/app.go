// Package app содержит инициализацию и запуск сервиса, согласно конфигурации
package app

import (
	"context"

	"route256.ozon.ru/project/notifier/internal/providers/eventhandler"
	"route256.ozon.ru/project/notifier/internal/providers/handlingrunner"
	"route256.ozon.ru/project/notifier/internal/service"
)

// App это приложение, которое инициализирует и запускает notifier
type App struct {
	config  Config
	service *service.Notifier
}

// NewApp возращает App, готовый к запуску
func NewApp(config Config) *App {
	app := &App{
		config: config,
	}
	app.init()
	return app
}

func (a *App) init() {
	runner := handlingrunner.NewKafkaConsumerGroupRunner(a.config.Brokers, a.config.Topic)
	handler := eventhandler.NewLoggerToWriter(a.config.LogsWriter)
	a.service = service.NewNotifier(handler, runner)
}

// Run запускает notifier. Блокирующий.
func (a *App) Run(ctx context.Context) error {
	return a.service.Run(ctx)
}
