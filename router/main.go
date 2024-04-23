package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func MainRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/register", register)
	r.Get("/login", login)

	r.Route("/users", func(r chi.Router) {
		r.Get("/", getAllUser) // POST /users
		r.Post("/", addUser)
		r.Delete("/{userID}", removeUser)
	})
	return r
}
