package main

import (
	"context"
	"go.uber.org/zap"
	"os/signal"
	"route256.ozon.ru/project/cart/internal/app"
	"syscall"
)

func main() {
	logger := getLogger()
	defer logger.Sync()
	config, err := setupCartConfig()
	if err != nil {
		logger.Error("error setting up cart config", zap.Error(err))
		return
	}
	config.Logger = logger.With(zap.String("service", "cart"))
	appl, err := app.NewApp(config)
	if err != nil {
		logger.Error("error initializing cart app", zap.Error(err))
		return
	}
	processLiveContext, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	shutdown, err := setUpProductionTracing()
	if err != nil {
		logger.Error("error setting up tracing", zap.Error(err))
		return
	}
	defer func() {
		err := shutdown()
		if err != nil {
			logger.Error("error shutting down tracing", zap.Error(err))
		}
	}()

	err = appl.Run(processLiveContext)
	if err != nil {
		logger.Error("services completed with error", zap.Error(err))
	}
}
