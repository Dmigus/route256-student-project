package main

import (
	"log"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getLogger() *zap.Logger {
	cfg := zap.NewProductionConfig()
	level, err := zap.ParseAtomicLevel(viper.GetString("LoggerLevel"))
	if err != nil {
		log.Fatal(err)
	}
	cfg.Level = level
	cfg.OutputPaths = []string{"stdout"}
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	logger := zap.Must(cfg.Build(zap.WithCaller(false)))
	return logger
}
