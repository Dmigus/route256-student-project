package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"route256.ozon.ru/project/loms/internal/apps"
	"route256.ozon.ru/project/loms/internal/apps/loms"
	"route256.ozon.ru/project/loms/internal/apps/outboxsender"
)

func setupLOMSConfig() (loms.Config, error) {
	config, err := apps.NewConfig[loms.Config](cliFlags.lomsConfig)
	if err != nil {
		return config, err
	}
	if config.Storages == nil {
		return config, fmt.Errorf("storages must be set")
	}

	dbPassFromEnv := os.Getenv("POSTGRES_PASSWORD_FILE")
	if len(dbPassFromEnv) > 0 {
		postgresPwd, err := readSecretFromFile(dbPassFromEnv)
		if err != nil {
			return loms.Config{}, err
		}
		for i := range config.Storages {
			config.Storages[i].Master.Password = postgresPwd
			config.Storages[i].Replica.Password = postgresPwd
		}
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
		for i := range config.Outbox {
			config.Outbox[i].Password = postgresPwd
		}
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
