package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/mariaulfa33/go-auth/usecase"
)

func MainRouter(userController usecase.Users) http.Handler {
	r := chi.NewRouter()
	r.Post("/register", userController.Register)
	r.Post("/login", userController.Login)

	r.Route("/users", func(r chi.Router) {
		r.Get("/", userController.GetAllUser)
		r.Post("/", userController.AddUser)
		r.Delete("/{userID}", userController.RemoveUserById)
	})
	return r
}
