// Package main содержит main для запуска программы, а также обработку переменных окружения и аргументов командной строки
package main

import (
	"context"
	_ "net/http/pprof"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"
	"route256.ozon.ru/project/loms/internal/apps/loms"
	"route256.ozon.ru/project/loms/internal/apps/outboxsender"
)

func main() {
	logger := getLogger()
	defer logger.Sync()

	lomsConfig, err := setupLOMSConfig()
	if err != nil {
		logger.Error("error setting up loms config", zap.Error(err))
		return
	}
	lomsConfig.Logger = logger.With(zap.String("service", "loms"))
	lomsApp, err := loms.NewApp(lomsConfig)
	if err != nil {
		logger.Error("error initializing loms app", zap.Error(err))
		return
	}
	outboxSenderConfig, err := setupOutboxSenderConfig()
	if err != nil {
		logger.Error("error setting up outbox sender config", zap.Error(err))
		return
	}
	outboxSenderConfig.Logger = logger.With(zap.String("service", "outbox_sender"))
	outboxSenderApp, err := outboxsender.NewApp(outboxSenderConfig)
	if err != nil {
		logger.Error("error initializing outbox sender app", zap.Error(err))
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
		errShutdown := shutdown()
		if errShutdown != nil {
			logger.Error("error shutting down tracing", zap.Error(errShutdown))
		}
	}()

	wg := &sync.WaitGroup{}
	wg.Add(2)
	var lomsErr error
	go func() {
		defer wg.Done()
		defer stop()
		lomsErr = lomsApp.Run(processLiveContext)
	}()
	go func() {
		defer wg.Done()
		defer stop()
		outboxSenderApp.Run(processLiveContext)
	}()
	wg.Wait()
	if lomsErr != nil {
		logger.Error("services completed with error", zap.Error(err))
	}
}
