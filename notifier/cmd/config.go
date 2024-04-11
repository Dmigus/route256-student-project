package main

import (
	"route256.ozon.ru/project/notifier/internal/app"
)

func setupConfig() (app.Config, error) {
	return app.NewConfig[app.Config](cliFlags.config)
}
