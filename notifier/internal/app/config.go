package app

import (
	"go.uber.org/zap"
)

// Config это конфигурация приложения
type Config struct {
	Brokers []string `json:"Brokers"`
	Topic   string   `json:"Topic"`
	Logger  *zap.Logger
}
