package handler

import (
	"net/http"

	"github.com/jackysum/go-template/web/template/page"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Root(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.URL.Path == "/" {
		w.WriteHeader(http.StatusOK)
		page.Home().Render(r.Context(), w)

		return
	}

	w.WriteHeader(http.StatusNotFound)
	page.NotFound().Render(r.Context(), w)
}
