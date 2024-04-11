package main

import "flag"

type flags struct {
	config string
}

var cliFlags = flags{}

func init() {
	flag.StringVar(&cliFlags.config, "config", "./configs/local.json", "path to config file for notifier")
	flag.Parse()
}
