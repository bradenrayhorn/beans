package inmem

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"

	"github.com/bradenrayhorn/beans/server/beans"
)

type sessionRepository struct {
	sessions map[string]beans.Session
	mu       sync.RWMutex
}

func NewSessionRepository() *sessionRepository {
	return &sessionRepository{
		sessions: make(map[string]beans.Session),
	}
}

func (r *sessionRepository) Create(userID beans.ID) (beans.Session, error) {
	bytes := make([]byte, 64)
	_, err := rand.Read(bytes)
	if err != nil {
		return beans.Session{}, err
	}
	sessionID := base64.RawURLEncoding.EncodeToString(bytes)

	// lock map
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.sessions[sessionID]; ok {
		return beans.Session{}, errors.New("session id conflict")
	}

	session := beans.Session{
		ID:     beans.SessionID(sessionID),
		UserID: userID,
	}

	r.sessions[sessionID] = session

	return session, nil
}

func (r *sessionRepository) Get(id beans.SessionID) (beans.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if session, ok := r.sessions[string(id)]; ok {
		return session, nil
	}

	return beans.Session{}, beans.ErrorNotFound
}

func (r *sessionRepository) Delete(id beans.SessionID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.sessions, string(id))

	return nil
}
