package session

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type SessionStore struct {
	db *sql.DB
}

func NewSessionStore(db *sql.DB) *SessionStore {
	return &SessionStore{db: db}
}

func (s *SessionStore) CreateSession(ctx context.Context, userID string) (Session, error) {
	session := Session{
		ID:         uuid.NewString(),
		UserID:     userID,
		Expiration: time.Now().Add(time.Minute * 30),
	}

	_, err := s.db.Exec("INSERT INTO session (id, userId, expiration) VALUES (?, ?, ?);", session.ID, session.UserID, session.Expiration)
	if err != nil {
		return Session{}, fmt.Errorf("failed to create session for user: %s \n Error: %w", userID, err)
	}

	return session, nil
}
