package app

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

type Config struct {
	Server struct {
		Host                  string `json:"Host"`
		Port                  uint16 `json:"Port"`
		ShutdownTimoutSeconds uint   `json:"ShutdownTimoutSeconds"`
	} `json:"Server"`
	ProductService struct {
		BaseURL     string `json:"BaseURL"`
		AccessToken string `json:"AccessToken"`
		RetryPolicy struct {
			RetryStatusCodes []int `json:"RetryStatusCodes"`
			MaxRetries       int   `json:"MaxRetries"`
		} `json:"RetryPolicy"`
		RPS uint64 `json:"RPS"`
	} `json:"ProductService"`
	LOMS struct {
		Address string `json:"Address"`
	}
	Logger *zap.Logger
	Redis  struct {
		Addr string `json:"Addr"`
	} `json:"Redis"`
}

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
