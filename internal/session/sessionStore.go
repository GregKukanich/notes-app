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
		Expiration: time.Now().Add(time.Minute * 30).Unix(),
	}

	_, err := s.db.Exec("INSERT INTO session (id, userId, expiration) VALUES (?, ?, ?);", session.ID, session.UserID, session.Expiration)
	if err != nil {
		return Session{}, fmt.Errorf("failed to create session for user: %s \n Error: %w", userID, err)
	}

	return session, nil
}

func (s *SessionStore) CheckSession(ctx context.Context, sessionId string) (bool, error) {
	var session Session
	row := s.db.QueryRow("SELECT id, userId, expiration FROM session WHERE id=?;", sessionId)
	err := row.Scan(&session.ID, &session.UserID, &session.Expiration)
	if err != nil {
		return false, fmt.Errorf("failed to find session matching sessionId %s: %w", sessionId, err)
	}

	cur := time.Now().Unix()
	if session.Expiration <= cur {
		return false, fmt.Errorf("session: %s is expired", sessionId)
	}

	return true, nil

}
