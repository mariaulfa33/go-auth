package usecase

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mariaulfa33/go-auth/repository"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	UserService *repository.UserService
}

type UserAuthRequest struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type UserResponse struct {
	User *repository.UserResponse `json:"user"`
}

type UserLoginResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

// Render implements render.Renderer.
func (u *UserLoginResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type ErrResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	AppCode    int64  `json:"code,omitempty"`
	ErrorText  string `json:"error,omitempty"`
}

var Secret = []byte(os.Getenv("JWT_SECRET"))

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrServerError(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Something unexpected error",
		ErrorText:      err.Error(),
	}
}
func (a *UserAuthRequest) Bind(r *http.Request) error {
	if a.Email == nil || a.Password == nil {
		return errors.New("missing required Register fields")
	}

	return nil
}

func (u *UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func UserResponseData(user *repository.UserResponse) render.Renderer {
	resp := UserResponse{user}

	return &resp
}

func UserResponseDataWithToken(user UserLoginResponse) render.Renderer {
	resp := UserLoginResponse{
		Id:    user.Id,
		Email: user.Email,
		Token: user.Token,
	}

	return &resp
}

func hashPassword(password string) string {
	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}

	return string(hashedBytes)
}

func userListResponse(users []repository.UserResponse) []render.Renderer {
	list := []render.Renderer{}

	//@TODO: Refactor this after learning more about pointer
	for _, user := range users {
		list = append(list, UserResponseData(&repository.UserResponse{
			Id:    user.Id,
			Email: user.Email,
		}))
	}

	return list
}

func (u Users) Register(w http.ResponseWriter, r *http.Request) {
	data := &UserAuthRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	userRequest := repository.CreateUser{
		Email:        strings.ToLower(*data.Email),
		PasswordHash: hashPassword(*data.Password),
	}

	existingUser, _ := u.UserService.GetUserByEmail(userRequest.Email)

	if existingUser != nil {
		render.Render(w, r, ErrInvalidRequest(errors.New("user already exist")))
		return
	}
	user, err := u.UserService.CreateNewUser(userRequest)

	if err != nil {
		render.Render(w, r, ErrServerError(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, UserResponseData(&repository.UserResponse{
		Id:    user.Id,
		Email: user.Email,
	}))
}

func createJWTToken(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":  id,
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(Secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (u Users) Login(w http.ResponseWriter, r *http.Request) {
	data := &UserAuthRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	user, err := u.UserService.GetUserByEmail(*data.Email)

	if err == sql.ErrNoRows {
		render.Render(w, r, ErrInvalidRequest(errors.New("Users Not Found")))
		return
	}

	if err != nil {
		render.Render(w, r, ErrServerError(err))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(*data.Password))
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(errors.New("Users Not Found")))
		return
	}

	token, err := createJWTToken(user.Id)
	if err != nil {
		return
	}
	fmt.Println(token)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, UserResponseDataWithToken(UserLoginResponse{
		Id:    user.Id,
		Email: user.Email,
		Token: token,
	}))
}

func (u Users) GetAllUser(w http.ResponseWriter, r *http.Request) {
	users, err := u.UserService.GetAllUser()

	if err != nil {
		render.Render(w, r, ErrServerError(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.RenderList(w, r, userListResponse(users))
}

func (u Users) AddUser(w http.ResponseWriter, r *http.Request) {
	data := &UserAuthRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	userRequest := repository.CreateUser{
		Email:        strings.ToLower(*data.Email),
		PasswordHash: hashPassword(*data.Password),
	}

	existingUser, _ := u.UserService.GetUserByEmail(userRequest.Email)

	if existingUser != nil {
		render.Render(w, r, ErrInvalidRequest(errors.New("user already exist")))
		return
	}

	user, err := u.UserService.CreateNewUser(userRequest)

	if err != nil {
		render.Render(w, r, ErrServerError(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, UserResponseData(&repository.UserResponse{
		Id:    user.Id,
		Email: user.Email,
	}))

}

func (u Users) RemoveUserById(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	if userID == "" {
		render.Render(w, r, ErrInvalidRequest(errors.New("Users Not Found")))
	}

	_, err := u.UserService.GetUserById(userID)

	if err == sql.ErrNoRows {
		render.Render(w, r, ErrInvalidRequest(errors.New("Users Not Found")))
		return
	}

	if err != nil {
		render.Render(w, r, ErrServerError(err))
		return
	}

	deletedUser, err := u.UserService.DeleteUser(userID)

	if err != nil {
		render.Render(w, r, ErrServerError(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, UserResponseData(&repository.UserResponse{
		Email: deletedUser.Email,
		Id:    deletedUser.Id,
	}))

}
