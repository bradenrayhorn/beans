package inmem_test

import (
	"sync"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/inmem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCanCreateAndGetSession(t *testing.T) {
	r := inmem.NewSessionRepository()

	userID := beans.NewBeansID()
	session, err := r.Create(userID)
	require.Nil(t, err)

	assert.Equal(t, userID, session.UserID)
	assert.Greater(t, len(session.ID), 0)

	gotSession, err := r.Get(session.ID)
	require.Nil(t, err)
	assert.Equal(t, session.ID, gotSession.ID)
	assert.Equal(t, session.UserID, gotSession.UserID)
}

func TestCannotGetNonExistent(t *testing.T) {
	r := inmem.NewSessionRepository()

	_, err := r.Get(beans.SessionID("blah"))
	assert.ErrorIs(t, err, beans.ErrorNotFound)
}

func TestCanDeleteSession(t *testing.T) {
	r := inmem.NewSessionRepository()

	userID := beans.NewBeansID()
	session, err := r.Create(userID)
	require.Nil(t, err)

	_, err = r.Get(session.ID)
	require.Nil(t, err)

	err = r.Delete(session.ID)
	require.Nil(t, err)

	_, err = r.Get(session.ID)
	require.Equal(t, beans.ErrorNotFound, err)
}

func TestCanUseConcurrentSessions(t *testing.T) {
	r := inmem.NewSessionRepository()
	var wg sync.WaitGroup

	makeAndGetSession := func(i int) {
		defer wg.Done()
		userID := beans.NewBeansID()
		session, err := r.Create(userID)
		require.Nil(t, err)

		gotSession, err := r.Get(session.ID)
		require.Nil(t, err)
		assert.Equal(t, session.ID, gotSession.ID)
		assert.Equal(t, session.UserID, gotSession.UserID)

		require.Nil(t, r.Delete(session.ID))
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go makeAndGetSession(i)
	}

	wg.Wait()
}
