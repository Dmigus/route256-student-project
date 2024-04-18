package main

import "flag"

type flags struct {
	lomsConfig         string
	outboxSenderConfig string
	// OTEL_EXPORTER_OTLP_TRACES_ENDPOINT = http://127.0.0.1:4318/v1/traces
}

var cliFlags = flags{}

func init() {
	flag.StringVar(&cliFlags.lomsConfig, "loms-config", "./configs/loms/local.json", "path to config file for loms")
	flag.StringVar(&cliFlags.outboxSenderConfig, "outbox-sender-config", "./configs/outbox-sender/local.json", "path to config file for outbox sender")
	flag.Parse()
}
