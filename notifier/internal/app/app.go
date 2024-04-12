// Package app содержит инициализацию и запуск сервиса, согласно конфигурацию
package app

import (
	"context"
	"errors"
	"route256.ozon.ru/project/notifier/internal/providers/eventhandler"
	"route256.ozon.ru/project/notifier/internal/providers/handlingrunner"
	"route256.ozon.ru/project/notifier/internal/service"
	"strconv"
	"sync"
	"time"
)

// App это приложение, которое инициализирует и запускает экземпляры notifier'ов
type App struct {
	config   Config
	services []*service.Notifier
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
	for i := 0; i < a.config.NotifiersNum; i++ {
		handlerID := strconv.Itoa(i)
		handler := eventhandler.NewLoggerToWriter(a.config.LogsWriter, handlerID)
		a.services = append(a.services, service.NewNotifier(handler, runner))
	}
}

// Run запускает сконфигурированное количество notifier'ов. Блокирующий. Завершится, когда завершаться все notifier'ы
func (a *App) Run(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	errs := make([]error, len(a.services))
	for ind, notifier := range a.services {
		wg.Add(1)
		ind := ind
		notifier := notifier
		go func() {
			defer wg.Done()
			errs[ind] = notifier.Run(ctx)
		}()
		time.Sleep(5 * time.Second)
	}
	wg.Wait()
	return errors.Join(errs...)
}
