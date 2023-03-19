package inmem

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"

	"github.com/bradenrayhorn/beans/beans"
)

type sessionRepository struct {
	sessions map[string]*beans.Session
	mu       sync.RWMutex
}

func NewSessionRepository() *sessionRepository {
	return &sessionRepository{
		sessions: make(map[string]*beans.Session),
	}
}

func (r *sessionRepository) Create(userID beans.ID) (*beans.Session, error) {
	bytes := make([]byte, 64)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	sessionID := base64.RawURLEncoding.EncodeToString(bytes)

	// lock map
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.sessions[sessionID]; ok {
		return nil, beans.WrapError(errors.New("session id conflict"), beans.ErrorInternal)
	}

	session := &beans.Session{
		ID:        beans.SessionID(sessionID),
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	r.sessions[sessionID] = session

	return session, nil
}

func (r *sessionRepository) Get(id beans.SessionID) (*beans.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if session, ok := r.sessions[string(id)]; ok {
		return session, nil
	}

	return nil, beans.ErrorNotFound
}

func (r *sessionRepository) Delete(id beans.SessionID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.sessions, string(id))

	return nil
}
