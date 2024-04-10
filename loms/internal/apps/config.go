// Package apps содержит единицы, которые могут выполнять некоторую логику.
package apps

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// NewConfig читает файл configPath в формате json в структуру типа configType и возвращает её
func NewConfig[configType any](configPath string) (configType, error) {
	var conf configType
	file, err := os.Open(filepath.Clean(configPath))
	if err != nil {
		return conf, err
	}
	defer func() {
		err2 := file.Close()
		err = errors.Join(err, err2)
	}()
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&conf); err != nil {
		return conf, err
	}
	return conf, nil
}

// PostgresConnectConfig это конфиг для подключения к PostgreSQL
type PostgresConnectConfig struct {
	User     string `json:"User"`
	Host     string `json:"Host"`
	Port     uint16 `json:"Port"`
	Database string `json:"Database"`
	Password string `json:"Password"`
}

// GetPostgresDSN возвращает Data Source Name, согласно конфигурации
func (pc PostgresConnectConfig) GetPostgresDSN() string {
	host := fmt.Sprintf("host=%s", pc.Host)
	user := fmt.Sprintf("user=%s", pc.User)
	password := fmt.Sprintf("password=%s", pc.Password)
	port := fmt.Sprintf("port=%d", pc.Port)
	db := fmt.Sprintf("dbname=%s", pc.Database)
	sslmode := fmt.Sprintf("sslmode=%s", "disable")
	return strings.Join([]string{host, user, password, port, db, sslmode}, " ")
}
