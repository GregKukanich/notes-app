package main

import "context"

type User struct {
	ID       string
	UserName string
	Password string
}

type UserRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID        string `json:"id"`
	UserName  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

type UserStore interface {
	CreateUser(ctx context.Context, u User) (UserResponse, error)
	// GetByUserName(ctx context.Context, username string) (UserResponse, error)
	GetByID(ctx context.Context, id string) (UserResponse, error)
}
