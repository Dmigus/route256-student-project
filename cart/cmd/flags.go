package main

import "flag"

type flags struct {
	configPath  string
	loggerLevel string
}

var cliFlags = flags{}

func init() {
	flag.StringVar(&cliFlags.configPath, "config", "./configs/local.json", "path to config file")
	flag.StringVar(&cliFlags.loggerLevel, "logger-level", "info", "logger level")
	flag.Parse()
}
