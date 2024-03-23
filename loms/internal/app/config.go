package app

import (
	_ "embed"
	"encoding/json"
	"errors"
	"os"
)

//go:embed stock-data.json
var stockdata []byte

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
