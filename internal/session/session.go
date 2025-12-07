package session

import (
	"context"
)

type Session struct {
	ID         string
	UserID     string
	Expiration int64
}

type SessionStoreIntf interface {
	CreateSession(ctx context.Context, userID string) (Session, error)
}
