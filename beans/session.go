package beans

import "time"

type SessionID string

type Session struct {
	ID        SessionID
	UserID    UserID
	CreatedAt time.Time
}

type SessionRepository interface {
	Create(userID UserID) (*Session, error)
	Get(id SessionID) (*Session, error)
}
