package server

import (
	"net/http"

	"github.com/jackysum/go-template/src/server/handler"
	"github.com/jackysum/go-template/src/server/middleware"
	"github.com/rs/zerolog"
)

type middlewareOpt func(http.Handler) http.Handler

func Setup(h *handler.Handler, opts ...middlewareOpt) http.Handler {
	mux := http.NewServeMux()
	routes(mux, h)

	var s http.Handler = mux
	for _, opt := range opts {
		s = opt(s)
	}

	return s
}

func WithLogger(log zerolog.Logger) middlewareOpt {
	return func(next http.Handler) http.Handler {
		return middleware.NewLogger(next, log)
	}
}
