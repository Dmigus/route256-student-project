// Package main содержит main для запуска программы, а также обработку переменных окружения и аргументов командной строки
package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"route256.ozon.ru/project/notifier/internal/app"
)

func main() {
	config, err := setupConfig()
	if err != nil {
		log.Fatal(err)
	}
	notifierApp := app.NewApp(config)
	processLiveContext, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	shutdownProvider, err := setUpProductionTracing()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := shutdownProvider()
		if err != nil {
			log.Println(err)
		}
	}()

	err = notifierApp.Run(processLiveContext)
	if err != nil {
		log.Printf("%v\n", err)
	}
}
