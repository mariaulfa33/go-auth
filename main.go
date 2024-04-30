package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mariaulfa33/go-auth/repository"
	router "github.com/mariaulfa33/go-auth/routers"
	"github.com/mariaulfa33/go-auth/usecase"
)

func main() {
	//load env file
	usecase.LoadEnv(".env")

	// Setup the database
	db, err := repository.Open(repository.ReadPostgresConfig())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Set up services
	userService := &repository.UserService{
		DB: db,
	}

	// Set up controllers
	usersC := usecase.Users{
		UserService: userService,
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/", router.MainRouter(usersC))

	fmt.Println("Listening on port 3000...")
	http.ListenAndServe(":3000", r)
}
