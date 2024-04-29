package repository

import (
	"database/sql"
)

type User struct {
	Id           string
	Email        string
	Username     string
	PasswordHash string
}

type CreateUser struct {
	Username     string
	Email        string
	PasswordHash string
}

type UserResponse struct {
	Id       string
	Email    string
	Username string
}

type UserService struct {
	DB *sql.DB
}

func (us *UserService) CreateNewUser(user CreateUser) (*User, error) {
	row := us.DB.QueryRow(`
		INSERT INTO users (email, username, passwordHash)
		VALUES ($1, $2, $3) RETURNING id, email, username`,
		user.Email, user.Username, user.PasswordHash)
	var userResult User
	err := row.Scan(&userResult.Id, &userResult.Email, &userResult.Username)
	if err != nil {
		return nil, err
	}
	return &userResult, nil
}

func (us *UserService) GetUserByEmailAndUsername(email string, username string) (*User, error) {
	row := us.DB.QueryRow(`
		SELECT id, email, username, passwordHash FROM users WHERE email = ($1) OR username = ($2)`, email, username)

	var userResult User
	err := row.Scan(&userResult.Id, &userResult.Email,
		&userResult.Username, &userResult.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &userResult, nil
}

func (us *UserService) GetUserByUsername(username string) (*User, error) {
	row := us.DB.QueryRow(`
		SELECT id, email, username, passwordHash FROM users WHERE username = ($1)`, username)

	var userResult User
	err := row.Scan(&userResult.Id, &userResult.Email,
		&userResult.Username, &userResult.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &userResult, nil
}

func (us *UserService) GetUserById(id string) (*User, error) {
	row := us.DB.QueryRow(`
		SELECT * FROM users WHERE id = ($1)`, id)

	var userResult User
	err := row.Scan(&userResult.Id, &userResult.Email, &userResult.Username, &userResult.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &userResult, nil
}

func (us *UserService) GetAllUser() ([]UserResponse, error) {
	rows, err := us.DB.Query(`
		SELECT id, email, username FROM users`)

	if err != nil {
		return nil, err
	}

	var users []UserResponse

	for rows.Next() {
		var user UserResponse
		if err := rows.Scan(&user.Id, &user.Email, &user.Username); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return users, err
	}

	return users, nil

}

func (us *UserService) DeleteUser(id string) (*User, error) {
	row := us.DB.QueryRow(`DELETE FROM users WHERE id=($1) RETURNING *`, id)

	var userResult User
	err := row.Scan(&userResult.Id, &userResult.Email,
		&userResult.Username, &userResult.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &userResult, nil
}
