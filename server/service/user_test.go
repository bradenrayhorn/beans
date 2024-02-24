package service_test

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser(t *testing.T) {
	services, _, _, sessionRepository, stop := makeServices(t)
	defer stop()
	ctx := context.Background()

	t.Run("GetAuth", func(t *testing.T) {

		t.Run("can get", func(t *testing.T) {
			userID := beans.NewID()

			// make session
			session, err := sessionRepository.Create(userID)
			require.NoError(t, err)

			// get auth context and verify
			auth, err := services.User.GetAuth(ctx, session.ID)
			require.NoError(t, err)

			assert.Equal(t, session.ID, auth.SessionID())
			assert.Equal(t, userID, auth.UserID())
		})

		t.Run("gives unauthorized error with bad session", func(t *testing.T) {
			userID := beans.NewID()

			// make session
			_, err := sessionRepository.Create(userID)
			require.NoError(t, err)

			// get auth context with bogus session
			_, err = services.User.GetAuth(ctx, beans.SessionID("123"))
			testutils.AssertErrorCode(t, err, beans.EUNAUTHORIZED)
		})
	})

}
