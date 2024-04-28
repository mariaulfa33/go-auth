package usecase

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
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

type ErrResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	AppCode    int64  `json:"code,omitempty"`
	ErrorText  string `json:"error,omitempty"`
}

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
	for _, user := range users {
		list = append(list, getAllUserResponse(user)) // error
	}
	return list
}

func getAllUserResponse(user *repository.UserResponse) UserResponse {
	return &UserResponse{User: user} // error

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

	render.Status(r, http.StatusCreated)
	render.Render(w, r, UserResponseData(&repository.UserResponse{
		Id:    user.Id,
		Email: user.Email,
	}))
}

func (u Users) GetAllUser(w http.ResponseWriter, r *http.Request) {
	users, err := u.UserService.GetAllUser()

	if err != nil {
		render.Render(w, r, ErrServerError(err))
		return
	}

	fmt.Println(users, "users")
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
