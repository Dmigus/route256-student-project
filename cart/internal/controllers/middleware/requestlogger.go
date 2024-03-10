package middleware

import (
	"log"
	"net/http"
	"os"
)

type RequestLoggerMW struct {
	wrapped http.Handler
	logger  *log.Logger
}

func NewLogger(handlerToWrap http.Handler) *RequestLoggerMW {
	return &RequestLoggerMW{
		wrapped: handlerToWrap,
		logger:  log.New(os.Stdout, "Request received: ", log.Lmsgprefix|log.LstdFlags),
	}
}

func (rl *RequestLoggerMW) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rl.logger.Printf("%s %s\n", r.Method, r.URL.Path)
	rl.wrapped.ServeHTTP(w, r)
}
