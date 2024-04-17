package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"os"
	"path/filepath"

	"route256.ozon.ru/project/loms/internal/apps"
	"route256.ozon.ru/project/loms/internal/apps/loms"
	"route256.ozon.ru/project/loms/internal/apps/outboxsender"
)

func setupLOMSConfig() (loms.Config, error) {
	config, err := apps.NewConfig[loms.Config](cliFlags.lomsConfig)
	if err != nil {
		return config, err
	}
	if config.Storage == nil {
		return config, fmt.Errorf("storage must be set")
	}

	dbPassFromEnv := os.Getenv("POSTGRES_PASSWORD_FILE")
	if len(dbPassFromEnv) > 0 {
		postgresPwd, err := readSecretFromFile(dbPassFromEnv)
		if err != nil {
			return loms.Config{}, err
		}
		config.Storage.Master.Password = postgresPwd
		config.Storage.Replica.Password = postgresPwd
	}
	config.MetricsRegisterer = prometheus.DefaultRegisterer
	config.MetricsHandler = promhttp.Handler()
	return config, nil
}

func setupOutboxSenderConfig() (outboxsender.Config, error) {
	config, err := apps.NewConfig[outboxsender.Config](cliFlags.outboxSenderConfig)
	if err != nil {
		return outboxsender.Config{}, err
	}

	dbPassFromEnv := os.Getenv("POSTGRES_PASSWORD_FILE")
	if len(dbPassFromEnv) > 0 {
		postgresPwd, err := readSecretFromFile(dbPassFromEnv)
		if err != nil {
			return config, err
		}
		config.Outbox.Password = postgresPwd
	}
	config.MetricsRegisterer = prometheus.DefaultRegisterer
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
