package server

import (
	"net/http"

	"github.com/jackysum/go-template/src/server/handler"
	"github.com/jackysum/go-template/src/utils/file"
)

func routes(mux *http.ServeMux, h *handler.Handler) {
	fs := http.FileServer(http.Dir(file.AbsolutePath("web/static")))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", h.Root)
}
