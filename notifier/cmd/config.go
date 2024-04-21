package main

import (
	"route256.ozon.ru/project/notifier/internal/app"
)

func setupConfig() (app.Config, error) {
	cfg, err := app.NewConfig[app.Config](cliFlags.config)
	return cfg, err
}
