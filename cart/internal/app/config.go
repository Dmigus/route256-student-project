package app

import (
	"encoding/json"
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

func NewConfig(configPath string) (Config, error) {
	config := Config{}
	file, err := os.Open(configPath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&config); err != nil {
		return Config{}, err
	}
	return config, nil
}
