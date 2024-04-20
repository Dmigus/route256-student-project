package main

import "flag"

type flags struct {
	lomsConfig         string
	outboxSenderConfig string
	loggerLevel        string
}

var cliFlags = flags{}

func init() {
	flag.StringVar(&cliFlags.lomsConfig, "loms-config", "./configs/loms/local.json", "path to config file for loms")
	flag.StringVar(&cliFlags.outboxSenderConfig, "outbox-sender-config", "./configs/outbox-sender/local.json", "path to config file for outbox sender")
	flag.StringVar(&cliFlags.loggerLevel, "logger-level", "info", "logger level")
	flag.Parse()
}
