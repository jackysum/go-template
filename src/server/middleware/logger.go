package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type logger struct {
	next http.Handler
	log  zerolog.Logger
}

func NewLogger(next http.Handler, log zerolog.Logger) *logger {
	return &logger{
		next: next,
		log:  log,
	}
}

func (l *logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	lw := &loggerResponseWriter{
		ResponseWriter: w,
	}

	defer func() {
		l.log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Dur("duration", time.Since(start)).
			Int("status", lw.statusCode).
			Msg("incoming request")
	}()

	l.next.ServeHTTP(lw, r)
}

type loggerResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (l *loggerResponseWriter) WriteHeader(statusCode int) {
	l.statusCode = statusCode
	l.ResponseWriter.WriteHeader(statusCode)
}
