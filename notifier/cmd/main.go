// Package main содержит main для запуска программы, а также обработку переменных окружения и аргументов командной строки
package main

import (
	"context"
	"go.uber.org/zap"
	"os/signal"
	"syscall"

	"route256.ozon.ru/project/notifier/internal/app"
)

func main() {
	logger := getLogger()
	defer logger.Sync()

	config, err := setupConfig()
	if err != nil {
		logger.Fatal("error setting up notifier config", zap.Error(err))
	}
	config.Logger = logger.With(zap.String("service", "notifier"))
	notifierApp := app.NewApp(config)
	processLiveContext, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	shutdownProvider, err := setUpProductionTracing()
	if err != nil {
		logger.Error("error setting up tracing", zap.Error(err))
		return
	}
	defer func() {
		err := shutdownProvider()
		if err != nil {
			logger.Error("error shutting down tracing", zap.Error(err))
		}
	}()

	err = notifierApp.Run(processLiveContext)
	if err != nil {
		logger.Error("service notifier completed with error", zap.Error(err))
	}
}
