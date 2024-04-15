package app

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// Config это конфигурация приложения
type Config struct {
	Brokers    []string `json:"Brokers"`
	Topic      string   `json:"Topic"`
	LogsWriter io.Writer
}

// NewConfig читает файл configPath в формате json в структуру типа configType и возвращает её
func NewConfig[configType any](configPath string) (configType, error) {
	var conf configType
	file, err := os.Open(filepath.Clean(configPath))
	if err != nil {
		return conf, err
	}
	defer func() {
		errClosing := file.Close()
		err = errors.Join(err, errClosing)
	}()
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&conf); err != nil {
		return conf, err
	}
	return conf, nil
}
