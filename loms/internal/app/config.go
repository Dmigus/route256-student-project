package app

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	GRPCServer struct {
		Port                  uint16 `json:"Port"`
		ShutdownTimoutSeconds uint   `json:"ShutdownTimoutSeconds"`
	} `json:"GRPCServer"`
	Stocks struct {
		InitData string `json:"InitData"`
	}
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
