package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os/signal"
	"route256.ozon.ru/project/cart/internal/app"
	"syscall"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./configs/local.json", "path to config file")
	flag.Parse()
	config, err := app.NewConfig(configPath)
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

	appl, err := app.NewApp(config)
	if err != nil {
		log.Printf("err initializing app: %v\n", err)
		return
	}
	err = appl.Run(processLiveContext)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
}
