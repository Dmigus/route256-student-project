package main

import "flag"

type flags struct {
	configPath string
}

var cliFlags = flags{}

func init() {
	flag.StringVar(&cliFlags.configPath, "config", "./configs/local.json", "path to config file")
	flag.Parse()
}
