package repository

import (
	"database/sql"
	"fmt"
)

type User struct {
	Id           string
	Email        string
	PasswordHash string
}

type CreateUser struct {
	Email        string
	PasswordHash string
}

type UserResponse struct {
	Id    string
	Email string
}

type UserService struct {
	DB *sql.DB
}

func (us *UserService) CreateNewUser(user CreateUser) (*User, error) {
	row := us.DB.QueryRow(`
		INSERT INTO users (email, passwordHash)
		VALUES ($1, $2) RETURNING id, email`, user.Email, user.PasswordHash)
	var userResult User
	err := row.Scan(&userResult.Id, &userResult.Email)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &userResult, nil
}

func (us *UserService) GetUserByEmail(email string) (*User, error) {
	row := us.DB.QueryRow(`
		SELECT * FROM users WHERE email = ($1)`, email)

	var userResult User
	err := row.Scan(&userResult.Id, &userResult.Email, &userResult.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &userResult, nil
}

func (us *UserService) GetUserById(id string) (*User, error) {
	row := us.DB.QueryRow(`
		SELECT * FROM users WHERE id = ($1)`, id)

	var userResult User
	err := row.Scan(&userResult.Id, &userResult.Email, &userResult.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &userResult, nil
}

func (us *UserService) GetAllUser() ([]UserResponse, error) {
	rows, err := us.DB.Query(`
		SELECT id, email FROM users`)

	if err != nil {
		return nil, err
	}

	var users []UserResponse

	for rows.Next() {
		var user UserResponse
		if err := rows.Scan(&user.Id, &user.Email); err != nil {
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
	err := row.Scan(&userResult.Id, &userResult.Email, &userResult.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &userResult, nil
}
