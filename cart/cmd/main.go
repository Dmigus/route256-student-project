package main

import (
	"context"
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

	appLiveContext, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	appl, err := app.NewApp(config)
	if err != nil {
		log.Printf("err initializing app: %v\n", err)
		return
	}
	err = appl.Run(appLiveContext)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
}
