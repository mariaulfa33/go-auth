package router

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/mariaulfa33/go-auth/usecase"
)

func MainRouter(userController usecase.Users) http.Handler {
	r := chi.NewRouter()
	r.Post("/register", userController.Register)
	r.Post("/login", userController.Login)

	r.Route("/users", func(r chi.Router) {
		r.Use(authenticationHandler)
		r.Get("/", userController.GetAllUser)
		r.Post("/", userController.Register)
		r.Delete("/{userID}", userController.RemoveUserById)
	})
	return r
}

type contextKey string

var tokenKey = contextKey("user")

func setUserSession(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, tokenKey, id)
}

func GetUserSession(ctx context.Context) (string, bool) {
	tokenStr, ok := ctx.Value(tokenKey).(string)
	return tokenStr, ok
}

// type JwtValueClaims struct {
// 	id  string
// 	exp int64
// 	jwt.RegisteredClaims
// }

func authenticationHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")
		if token == "" {
			render.Render(w, r, usecase.ErrInvalidRequest(errors.New("user unauthorized")))
			return
		}
		splitToken := strings.Split(token, "Bearer ")

		if len(splitToken) == 1 {
			render.Render(w, r, usecase.ErrInvalidRequest(errors.New("user unauthorized")))
			return
		}
		token = splitToken[1]

		jwtTokenParsed, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(usecase.Secret), nil
		})

		if err != nil {
			render.Render(w, r, usecase.ErrInvalidRequest(err))
			return
		}

		claims := jwtTokenParsed.Claims.(*jwt.RegisteredClaims)

		ctx := setUserSession(r.Context(), claims.Issuer)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
