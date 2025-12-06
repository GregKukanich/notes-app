package user

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

type UserStoreIntf interface {
	CreateUser(ctx context.Context, u User) (UserResponse, error)
	GetUserPasswordHash(ctx context.Context, username string) (string, error)
	GetByID(ctx context.Context, id string) (UserResponse, error)
	GetUserByUserName(ctx context.Context, username string) (User, error)
}
