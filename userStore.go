package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type userStore struct {
	db *sql.DB
}

func (us *userStore) CreateUser(ctx context.Context, u User) (UserResponse, error) {
	id := uuid.NewString()
	u.ID = id

	var ur UserResponse
	row := us.db.QueryRow("INSERT INTO user (id, username, password) VALUES (?, ?, ?) RETURNING id, username, created_at", u.ID, u.UserName, u.Password)
	err := row.Scan(&ur.ID, &ur.UserName, &ur.CreatedAt)
	if err != nil {
		return UserResponse{}, fmt.Errorf("failed to save User: %s \n Error: %w", u, err)
	}

	return ur, nil
}

// func (us *userStore) GetByUserName(ctx context.Context, username string) (UserResponse, error) {
// 	var ur UserResponse
// 	row := us.db.QueryRow("SELECT id, username, created_at FROM user WHERE username=%s", username)
// 	err := row.Scan(&ur.ID, &ur.UserName, &ur.CreatedAt)
// 	if err != nil {
// 		return UserResponse{}, fmt.Errorf("failed to find User with username: %s", username)
// 	}

// 	return ur, nil
// }

func (us *userStore) GetByID(ctx context.Context, id string) (UserResponse, error) {
	var ur UserResponse
	row := us.db.QueryRow("SELECT id, username, created_at FROM user WHERE id=?", id)
	err := row.Scan(&ur.ID, &ur.UserName, &ur.CreatedAt)
	if err != nil {
		return UserResponse{}, fmt.Errorf("failed to find User with username: %s", id)
	}

	return ur, nil
}
