package main

import (
	"flag"
	"log"
	"route256.ozon.ru/project/cart/internal/app"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./configs/config.json", "path to config file")
	flag.Parse()
	config, err := app.NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	app.NewApp(config).Run()
}
