package inmem_test

import (
	"testing"
	"time"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/inmem"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCanCreateAndGetSession(t *testing.T) {
	r := inmem.NewSessionRepository()

	userID := beans.UserID(ksuid.New())
	sessionTimestamp := time.Now()
	session, err := r.Create(userID)
	require.Nil(t, err)

	assert.Equal(t, userID, session.UserID)
	assert.GreaterOrEqual(t, session.CreatedAt, sessionTimestamp)
	assert.Greater(t, len(session.ID), 0)

	gotSession, err := r.Get(session.ID)
	require.Nil(t, err)
	assert.Equal(t, session.ID, gotSession.ID)
	assert.Equal(t, session.CreatedAt, gotSession.CreatedAt)
	assert.Equal(t, session.UserID, gotSession.UserID)
}

func TestCannotGetNonExistent(t *testing.T) {
	r := inmem.NewSessionRepository()

	session, err := r.Get(beans.SessionID("blah"))
	assert.Nil(t, session)
	assert.ErrorIs(t, err, beans.ErrorNotFound)
}
