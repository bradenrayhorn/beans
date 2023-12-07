package beans

import "time"

type SessionID string

type Session struct {
	ID        SessionID
	UserID    ID
	CreatedAt time.Time
}

type SessionRepository interface {
	Create(userID ID) (*Session, error)
	Get(id SessionID) (*Session, error)
	Delete(id SessionID) error
}
