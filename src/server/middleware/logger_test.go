package middleware_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackysum/go-template/src/server/middleware"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestLogger_ServeHTTP(t *testing.T) {
	buf := &bytes.Buffer{}
	log := zerolog.New(buf)

	mux := http.NewServeMux()
	h := middleware.NewLogger(mux, log)

	method := http.MethodGet
	path := "/url/path"

	r := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)

	got := buf.String()

	require.Contains(t, got, "\"level\":\"info\"")
	require.Contains(t, got, "\"method\":\"GET\"")
	require.Contains(t, got, fmt.Sprintf("\"path\":\"%s\"", path))
	require.Contains(t, got, "\"duration\":")
	require.Contains(t, got, "\"status\":404")
	require.Contains(t, got, "\"message\":\"incoming request\"")
}
