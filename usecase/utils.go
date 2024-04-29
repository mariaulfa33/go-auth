package usecase

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mariaulfa33/go-auth/repository"
	"golang.org/x/crypto/bcrypt"
)

var Secret = []byte(os.Getenv("JWT_SECRET"))

func (u *UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (u *UserLoginResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func UserResponseData(user *repository.UserResponse) render.Renderer {
	resp := UserResponse{
		Id:       user.Id,
		Email:    user.Email,
		Username: user.Username,
	}

	return &resp
}

func UserResponseDataWithToken(user UserLoginResponse) render.Renderer {
	resp := UserLoginResponse{
		Id:       user.Id,
		Email:    user.Email,
		Username: user.Username,
		Token:    user.Token,
	}

	return &resp
}

func userListResponse(users []repository.UserResponse) []render.Renderer {
	list := []render.Renderer{}

	for _, user := range users {
		list = append(list, UserResponseData(&user))
	}

	return list
}

func createJWTToken(id string) (string, error) {

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Minute)),
		Issuer:    id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(Secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

func (a *UserAuthRequest) Bind(r *http.Request) error {
	if a.Email == nil || a.Password == nil || a.Username == nil {
		return errors.New("missing required Register fields")
	}

	if len(*a.Email) < 1 || len(*a.Password) < 1 || len(*a.Username) < 1 {
		return errors.New("missing required Register fields")
	}

	return nil
}

func (a *UserLoginRequest) Bind(r *http.Request) error {
	if a.Password == nil || a.Username == nil {
		return errors.New("missing required Register fields")
	}

	return nil
}
