package main

import (
	"route256.ozon.ru/project/cart/internal/app"
)

func setupCartConfig() (app.Config, error) {
	return app.NewConfig[app.Config](cliFlags.configPath)
}
