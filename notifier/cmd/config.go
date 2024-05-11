package main

import (
	"flag"
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"route256.ozon.ru/project/notifier/internal/app"
)

const configNameFlag = "config"

func init() {
	flag.String(configNameFlag, "./configs/local.json", "path to config file for notifier")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		panic(fmt.Errorf("fatal error binding flags: %w", err))
	}
	configName := viper.GetString(configNameFlag)
	viper.SetConfigFile(configName)
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func getNotifierConfig() (app.Config, error) {
	conf := app.Config{}
	err := viper.Unmarshal(&conf)
	if err != nil {
		return app.Config{}, fmt.Errorf("fatal error config file: %w", err)
	}
	return conf, nil
}
