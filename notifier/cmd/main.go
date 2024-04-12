// Package main содержит main для запуска программы, а также обработку переменных окружения и аргументов командной строки
package main

import (
	"context"
	"log"
	"os/signal"
	"route256.ozon.ru/project/notifier/internal/app"
	"syscall"
)

func main() {
	config, err := setupConfig()
	if err != nil {
		log.Fatal(err)
	}
	notifierApp := app.NewApp(config)
	processLiveContext, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err = notifierApp.Run(processLiveContext)
	if err != nil {
		log.Printf("%v\n", err)
	}
}
