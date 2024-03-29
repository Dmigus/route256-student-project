package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"route256.ozon.ru/project/loms/internal/app"
	"syscall"
)

func main() {
	config, err := setupConfig()
	if err != nil {
		log.Fatal(err)
	}
	appl := app.NewApp(config)
	defer appl.Stop()
	go appl.Run()
	defer appl.StopGateway()
	go appl.RunGateway()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}

func setupConfig() (app.Config, error) {
	var configPath string
	flag.StringVar(&configPath, "config", "./configs/local.json", "path to config file")
	flag.Parse()
	config, err := app.NewConfig(configPath)
	if err != nil {
		return app.Config{}, err
	}
	if config.Storage == nil {
		return config, nil
	}

	dbPassFromEnv := os.Getenv("POSTGRES_PASSWORD_FILE")
	if len(dbPassFromEnv) > 0 {
		postgresPwd, err := readSecretFromFile(dbPassFromEnv)
		if err != nil {
			return app.Config{}, err
		}
		config.Storage.Password = postgresPwd
	}
	return config, nil
}

func readSecretFromFile(addr string) (string, error) {
	cleaned := filepath.Clean(addr)
	dataBytes, err := os.ReadFile(cleaned)
	if err != nil {
		return "", err
	}
	return string(dataBytes), nil
}
