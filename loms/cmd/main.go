// Package main содержит main для запуска программы, а также обработку переменных окружения и аргументов командной строки
package main

import (
	"context"
	"errors"
	"log"
	"os/signal"
	"route256.ozon.ru/project/loms/internal/apps/loms"
	"route256.ozon.ru/project/loms/internal/apps/outboxsender"
	"sync"
	"syscall"
)

func main() {
	lomsConfig, err := setupLOMSConfig()
	if err != nil {
		log.Fatal(err)
	}
	lomsApp, err := loms.NewApp(lomsConfig)
	if err != nil {
		log.Fatal(err)
	}
	outboxSenderConfig, err := setupOutboxSenderConfig()
	if err != nil {
		log.Fatal(err)
	}
	outboxSenderApp, err := outboxsender.NewApp(outboxSenderConfig)
	if err != nil {
		log.Fatal(err)
	}

	processLiveContext, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	tracerProvider, err := setUpProductionTracing()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		errFlush := tracerProvider.ForceFlush(context.Background())
		errShutdown := tracerProvider.Shutdown(context.Background())
		errs := errors.Join(errFlush, errShutdown)
		if errs != nil {
			log.Println(errs)
		}
	}()

	wg := &sync.WaitGroup{}
	wg.Add(2)
	var lomsErr error
	go func() {
		defer wg.Done()
		lomsErr = lomsApp.Run(processLiveContext)
	}()
	go func() {
		defer wg.Done()
		outboxSenderApp.Run(processLiveContext)
	}()
	wg.Wait()
	if lomsErr != nil {
		log.Fatal(err)
	}
}
