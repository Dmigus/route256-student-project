package outboxsender

import (
	"github.com/prometheus/client_golang/prometheus"
	"route256.ozon.ru/project/loms/internal/apps"
)

// Config это конфигурация для приложения outboxsender
type Config struct {
	Outbox apps.PostgresConnectConfig `json:"Outbox"`
	Kafka  struct {
		Brokers []string `json:"Brokers"`
		Topic   string   `json:"Topic"`
	} `json:"Kafka"`
	BatchSize         int32 `json:"BatchSize"`
	BatchInterval     int64 `json:"BatchInterval"`
	MetricsRegisterer prometheus.Registerer
}
