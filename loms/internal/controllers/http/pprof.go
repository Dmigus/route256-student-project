package http

import (
	"net/http"
	"net/http/pprof"
	"net/url"
)

const pprofBasePath = "/debug/pprof/"

func pprofHandler() (http.Handler, error) {
	mux := http.NewServeMux()
	mux.HandleFunc(pprofBasePath, pprof.Index)
	cmdLinePath, err := url.JoinPath(pprofBasePath, "cmdline")
	if err != nil {
		return nil, err
	}
	mux.HandleFunc(cmdLinePath, pprof.Cmdline)
	profilePath, err := url.JoinPath(pprofBasePath, "profile")
	if err != nil {
		return nil, err
	}
	mux.HandleFunc(profilePath, pprof.Profile)
	symbolPath, err := url.JoinPath(pprofBasePath, "symbol")
	if err != nil {
		return nil, err
	}
	mux.HandleFunc(symbolPath, pprof.Symbol)
	tracePath, err := url.JoinPath(pprofBasePath, "trace")
	if err != nil {
		return nil, err
	}
	mux.HandleFunc(tracePath, pprof.Trace)
	return mux, nil
}
