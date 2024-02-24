package datasource

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testUser(t *testing.T, ds beans.DataSource) {
	factory := testutils.NewFactory(t, ds)
	userRepository := ds.UserRepository()
	ctx := context.Background()

	t.Run("create", func(t *testing.T) {

		t.Run("cannot duplicate id", func(t *testing.T) {
			user := beans.User{
				ID:           beans.NewID(),
				Username:     beans.Username(beans.NewID().String()),
				PasswordHash: beans.PasswordHash("x"),
			}

			err := userRepository.Create(ctx, user.ID, user.Username, user.PasswordHash)
			require.NoError(t, err)

			// create again, should fail
			err = userRepository.Create(ctx, user.ID, beans.Username(beans.NewID().String()), user.PasswordHash)
			require.Error(t, err)
		})

		t.Run("cannot duplicate username", func(t *testing.T) {
			user := beans.User{
				ID:           beans.NewID(),
				Username:     beans.Username(beans.NewID().String()),
				PasswordHash: beans.PasswordHash("x"),
			}

			err := userRepository.Create(ctx, user.ID, user.Username, user.PasswordHash)
			require.NoError(t, err)

			// create again, should fail
			err = userRepository.Create(ctx, beans.NewID(), user.Username, user.PasswordHash)
			require.Error(t, err)
		})
	})

	t.Run("can create and get", func(t *testing.T) {
		user := beans.User{
			ID:           beans.NewID(),
			Username:     beans.Username(beans.NewID().String()),
			PasswordHash: beans.PasswordHash("x"),
		}
		err := userRepository.Create(ctx, user.ID, user.Username, user.PasswordHash)
		require.NoError(t, err)

		// get and verify
		res, err := userRepository.Get(ctx, user.ID)
		require.NoError(t, err)

		assert.Equal(t, user, res)
	})

	t.Run("get", func(t *testing.T) {

		t.Run("cannot get non-existent", func(t *testing.T) {
			factory.User(beans.User{})

			// get with a random ID
			_, err := userRepository.Get(ctx, beans.NewID())
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})
	})

	t.Run("get by username", func(t *testing.T) {

		t.Run("cannot get non-existent", func(t *testing.T) {
			factory.User(beans.User{})

			// get with a random username
			_, err := userRepository.GetByUsername(ctx, beans.Username(beans.NewID().String()))
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})

		t.Run("can get", func(t *testing.T) {
			user := factory.User(beans.User{})

			// get and verify results
			res, err := userRepository.GetByUsername(ctx, user.Username)
			require.NoError(t, err)

			assert.Equal(t, user, res)
		})
	})

	t.Run("exists", func(t *testing.T) {

		t.Run("existing user", func(t *testing.T) {
			user := factory.User(beans.User{})

			// user should exist
			res, err := userRepository.Exists(ctx, user.Username)
			require.NoError(t, err)

			assert.Equal(t, true, res)
		})

		t.Run("non-existent user", func(t *testing.T) {
			factory.User(beans.User{})

			// user should not exist
			res, err := userRepository.Exists(ctx, beans.Username(beans.NewID().String()))
			require.NoError(t, err)

			assert.Equal(t, false, res)
		})
	})
}
