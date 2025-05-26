package server

import (
	"net/http"

	"github.com/jackysum/go-template/src/utils/file"
)

func routes(mux *http.ServeMux) {
	fs := http.FileServer(http.Dir(file.AbsolutePath("web/static")))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
}
