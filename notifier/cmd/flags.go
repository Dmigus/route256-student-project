package main

import "flag"

type flags struct {
	config      string
	loggerLevel string
}

var cliFlags = flags{}

func init() {
	flag.StringVar(&cliFlags.config, "config", "./configs/local.json", "path to config file for notifier")
	flag.StringVar(&cliFlags.loggerLevel, "logger-level", "info", "logger level")
	flag.Parse()
}
