package server

import "net/http"

func Setup() http.Handler {
	mux := http.NewServeMux()

	return mux
}
