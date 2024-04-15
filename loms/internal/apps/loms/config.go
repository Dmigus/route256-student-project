package loms

import (
	_ "embed"
	"route256.ozon.ru/project/loms/internal/apps"
)

//go:embed stock-data.json
var stockdata []byte

// Config это конфигурация для приложения loms
type Config struct {
	GRPCServer struct {
		Port                  uint16 `json:"Port"`
		ShutdownTimoutSeconds uint   `json:"ShutdownTimoutSeconds"`
	} `json:"GRPCServer"`
	HTTPGateway struct {
		Port                  uint16 `json:"Port"`
		ShutdownTimoutSeconds uint   `json:"ShutdownTimoutSeconds"`
	} `json:"HTTPGateway"`
	Swagger struct {
		Path string `json:"Path"`
	} `json:"Swagger"`
	Storage *struct {
		Master  apps.PostgresConnectConfig `json:"Master"`
		Replica apps.PostgresConnectConfig `json:"Replica"`
	} `json:"Storage"`
}