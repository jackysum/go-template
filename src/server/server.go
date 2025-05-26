package server

import (
	"net/http"

	"github.com/jackysum/go-template/src/server/middleware"
	"github.com/rs/zerolog"
)

type middlewareOpt func(http.Handler) http.Handler

func Setup(opts ...middlewareOpt) http.Handler {
	mux := http.NewServeMux()
	routes(mux)

	var h http.Handler = mux
	for _, opt := range opts {
		h = opt(h)
	}

	return h
}

func WithLogger(log zerolog.Logger) middlewareOpt {
	return func(next http.Handler) http.Handler {
		return middleware.NewLogger(next, log)
	}
}
