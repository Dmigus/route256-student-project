package main

import (
	"os"
	"route256.ozon.ru/project/notifier/internal/app"
)

func setupConfig() (app.Config, error) {
	cfg, err := app.NewConfig[app.Config](cliFlags.config)
	cfg.LogsWriter = os.Stdout
	return cfg, err
}
