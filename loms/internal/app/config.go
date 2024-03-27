package app

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
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
	Storage *struct {
		User     string `json:"User"`
		Host     string `json:"Host"`
		Port     uint16 `json:"Port"`
		Database string `json:"Database"`
		Password string `json:"Password"`
	} `json:"Storage"`
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

func (c Config) getPostgresDSN() string {
	if c.Storage == nil {
		return ""
	}
	host := fmt.Sprintf("host=%s", c.Storage.Host)
	user := fmt.Sprintf("user=%s", c.Storage.User)
	password := fmt.Sprintf("password=%s", c.Storage.Password)
	port := fmt.Sprintf("port=%d", c.Storage.Port)
	db := fmt.Sprintf("dbname=%s", c.Storage.Database)
	sslmode := fmt.Sprintf("sslmode=%s", "disable")
	return strings.Join([]string{host, user, password, port, db, sslmode}, " ")
}
