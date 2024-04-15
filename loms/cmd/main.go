// Package main содержит main для запуска программы, а также обработку переменных окружения и аргументов командной строки
package main

import (
	"context"
	"log"
	"os/signal"
	"route256.ozon.ru/project/loms/internal/apps/loms"
	"route256.ozon.ru/project/loms/internal/apps/outboxsender"
	"syscall"
)

func main() {
	lomsConfig, err := setupLOMSConfig()
	if err != nil {
		log.Fatal(err)
	}
	lomsApp := loms.NewApp(lomsConfig)
	defer lomsApp.Stop()
	go lomsApp.Run()
	defer lomsApp.StopGateway()
	go lomsApp.RunGateway()

	processLiveContext, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	outboxSenderConfig, err := setupOutboxSenderConfig()
	if err != nil {
		log.Fatal(err)
	}
	outboxSenderApp, err := outboxsender.NewApp(outboxSenderConfig)
	if err != nil {
		log.Fatal(err)
	}
	outboxSenderApp.Run(processLiveContext)
}
