package session

import (
	"context"
	"time"
)

type Session struct {
	ID         string
	UserID     string
	Expiration time.Time
}

type SessionStoreIntf interface {
	CreateSession(ctx context.Context, userID string) (Session, error)
}
