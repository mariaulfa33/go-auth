package router

import (
	"fmt"
	"net/http"
)

func sampleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, r.URL.Path)
}

type Router struct{}

func (route Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/authentication":
		sampleHandler(w, r)
	case "/users":
		sampleHandler(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}
