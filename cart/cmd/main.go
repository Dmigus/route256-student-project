package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"route256.ozon.ru/project/cart/internal/app"
	"syscall"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./configs/config.json", "path to config file")
	flag.Parse()
	config, err := app.NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	appl := app.NewApp(config)
	defer appl.Stop()
	go appl.Run()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
