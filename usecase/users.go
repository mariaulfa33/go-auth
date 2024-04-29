package usecase

import (
	"database/sql"
	"errors"
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
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type UserLoginRequest struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}

type UserResponse struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type UserLoginResponse struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

func (u Users) Register(w http.ResponseWriter, r *http.Request) {
	data := &UserAuthRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	pass, err := hashPassword(*data.Password)
	if err != nil {
		render.Render(w, r, ErrServerError(err))
	}

	userRequest := repository.CreateUser{
		Username:     (*data.Username),
		Email:        strings.ToLower(*data.Email),
		PasswordHash: pass,
	}

	existingUser, _ := u.UserService.GetUserByEmailAndUsername(userRequest.Email, userRequest.Username)

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
		Id:       user.Id,
		Email:    user.Email,
		Username: user.Username,
	}))
}

func (u Users) Login(w http.ResponseWriter, r *http.Request) {
	data := &UserLoginRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	user, err := u.UserService.GetUserByUsername(*data.Username)

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
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	token, err := createJWTToken(user.Id)
	if err != nil {
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, UserResponseDataWithToken(UserLoginResponse{
		Id:       user.Id,
		Email:    user.Email,
		Username: user.Username,
		Token:    token,
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
		Id:       deletedUser.Id,
		Username: deletedUser.Username,
		Email:    deletedUser.Email,
	}))
}
