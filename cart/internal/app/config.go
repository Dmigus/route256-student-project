package app

import (
	"encoding/json"
	"errors"
	"os"
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
	} `json:"ProductService"`
}

func NewConfig(configPath string) (conf Config, err error) {
	file, err := os.Open(configPath)
	if err != nil {
		return Config{}, err
	}
	defer func() {
		err2 := file.Close()
		err = errors.Join(err, err2)
	}()
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&conf); err != nil {
		return Config{}, err
	}
	return conf, nil
}
