package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (us *UserStore) CreateUser(ctx context.Context, u User) (UserResponse, error) {
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

func (us *UserStore) GetUserPasswordHash(ctx context.Context, username string) (string, error) {
	var hash string
	row := us.db.QueryRow("SELECT password FROM user WHERE username=?", username)
	err := row.Scan(&hash)
	if err != nil {
		return "", fmt.Errorf("failed to find User with username: %s", username)
	}

	return hash, nil
}

func (us *UserStore) GetByID(ctx context.Context, id string) (UserResponse, error) {
	var ur UserResponse
	row := us.db.QueryRow("SELECT id, username, created_at FROM user WHERE id=?", id)
	err := row.Scan(&ur.ID, &ur.UserName, &ur.CreatedAt)
	if err != nil {
		return UserResponse{}, fmt.Errorf("failed to find User with id: %s", id)
	}

	return ur, nil
}

func (us *UserStore) GetUserByUserName(ctx context.Context, username string) (User, error) {
	var ur User
	row := us.db.QueryRow("SELECT id, username FROM user WHERE username=?", username)
	err := row.Scan(&ur.ID, &ur.UserName)
	if err != nil {
		return User{}, fmt.Errorf("failed to find User with username: %s", username)
	}

	return ur, nil
}
