package contract_test

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/contract"
	"github.com/bradenrayhorn/beans/server/inmem"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/bradenrayhorn/beans/server/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser(t *testing.T) {
	t.Parallel()
	pool, stop := testutils.StartPool(t)
	defer stop()

	cleanup := func() {
		_, err := pool.Exec(context.Background(), "truncate table users, budgets cascade;")
		require.Nil(t, err)
	}

	userRepository := postgres.NewUserRepository(pool)
	sessionRepository := inmem.NewSessionRepository()
	c := contract.NewUserContract(sessionRepository, userRepository)

	t.Run("register", func(t *testing.T) {
		t.Run("handles validation error", func(t *testing.T) {
			defer cleanup()

			err := c.Register(context.Background(), beans.Username(""), beans.Password(""))
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("can create user", func(t *testing.T) {
			defer cleanup()

			err := c.Register(context.Background(), beans.Username("user"), beans.Password("pass"))
			require.Nil(t, err)

			// user was saved
			_, err = userRepository.GetByUsername(context.Background(), beans.Username("user"))
			require.Nil(t, err)
		})

		t.Run("cannot create same user twice", func(t *testing.T) {
			defer cleanup()

			err := c.Register(context.Background(), beans.Username("user"), beans.Password("pass"))
			require.Nil(t, err)

			err = c.Register(context.Background(), beans.Username("user"), beans.Password("password"))
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})
	})

	t.Run("login", func(t *testing.T) {
		t.Run("username is required", func(t *testing.T) {
			defer cleanup()

			_, err := c.Login(context.Background(), beans.Username(""), beans.Password("pass"))
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("password is required", func(t *testing.T) {
			defer cleanup()

			_, err := c.Login(context.Background(), beans.Username("user"), beans.Password(""))
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("cannot login with non existant user", func(t *testing.T) {
			defer cleanup()

			_, err := c.Login(context.Background(), beans.Username("user"), beans.Password("pass"))
			testutils.AssertErrorAndCode(t, err, beans.EUNAUTHORIZED, "Invalid username or password")
		})

		t.Run("cannot login with invalid password", func(t *testing.T) {
			defer cleanup()

			err := c.Register(context.Background(), beans.Username("user"), beans.Password("pass"))
			require.Nil(t, err)

			_, err = c.Login(context.Background(), beans.Username("user"), beans.Password("password"))
			testutils.AssertErrorAndCode(t, err, beans.EUNAUTHORIZED, "Invalid username or password")
		})

		t.Run("can login", func(t *testing.T) {
			defer cleanup()

			err := c.Register(context.Background(), beans.Username("user"), beans.Password("pass"))
			require.Nil(t, err)

			session, err := c.Login(context.Background(), beans.Username("user"), beans.Password("pass"))
			require.Nil(t, err)

			dbUser, err := userRepository.Get(context.Background(), session.UserID)
			require.Nil(t, err)
			assert.Equal(t, "user", string(dbUser.Username))
		})
	})

	t.Run("logout", func(t *testing.T) {
		t.Run("can logout", func(t *testing.T) {
			defer cleanup()

			err := c.Register(context.Background(), beans.Username("user"), beans.Password("pass"))
			require.Nil(t, err)

			session, err := c.Login(context.Background(), beans.Username("user"), beans.Password("pass"))
			require.Nil(t, err)

			err = c.Logout(context.Background(), beans.NewAuthContext(session.UserID, session.ID))
			require.Nil(t, err)

			_, err = sessionRepository.Get(session.ID)
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})
	})

	t.Run("get me", func(t *testing.T) {
		t.Run("can get me", func(t *testing.T) {
			defer cleanup()

			err := c.Register(context.Background(), beans.Username("user"), beans.Password("pass"))
			require.Nil(t, err)

			session, err := c.Login(context.Background(), beans.Username("user"), beans.Password("pass"))
			require.Nil(t, err)

			user, err := c.GetMe(context.Background(), beans.NewAuthContext(session.UserID, session.ID))
			require.Nil(t, err)

			assert.Equal(t, "user", string(user.Username))
		})
	})
}
