package contract_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/contract"
	"github.com/bradenrayhorn/beans/internal/testutils"
	"github.com/bradenrayhorn/beans/postgres"
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
	c := contract.NewUserContract(userRepository)

	t.Run("create", func(t *testing.T) {
		t.Run("handles validation error", func(t *testing.T) {
			defer cleanup()

			_, err := c.CreateUser(context.Background(), beans.Username(""), beans.Password(""))
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("can create user", func(t *testing.T) {
			defer cleanup()

			user, err := c.CreateUser(context.Background(), beans.Username("user"), beans.Password("pass"))
			require.Nil(t, err)

			// user was returned
			assert.Equal(t, "user", string(user.Username))
			assert.False(t, user.ID.Empty())
			assert.NotEmpty(t, string(user.PasswordHash))

			// user was saved
			dbUser, err := userRepository.Get(context.Background(), user.ID)
			require.Nil(t, err)
			assert.True(t, reflect.DeepEqual(user, dbUser))
		})

		t.Run("cannot create same user twice", func(t *testing.T) {
			defer cleanup()

			_, err := c.CreateUser(context.Background(), beans.Username("user"), beans.Password("pass"))
			require.Nil(t, err)

			_, err = c.CreateUser(context.Background(), beans.Username("user"), beans.Password("password"))
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

			_, err := c.CreateUser(context.Background(), beans.Username("user"), beans.Password("pass"))
			require.Nil(t, err)

			_, err = c.Login(context.Background(), beans.Username("user"), beans.Password("password"))
			testutils.AssertErrorAndCode(t, err, beans.EUNAUTHORIZED, "Invalid username or password")
		})

		t.Run("can login", func(t *testing.T) {
			defer cleanup()

			user, err := c.CreateUser(context.Background(), beans.Username("user"), beans.Password("pass"))
			require.Nil(t, err)

			loggedInUser, err := c.Login(context.Background(), beans.Username("user"), beans.Password("pass"))
			require.Nil(t, err)

			assert.Equal(t, user.ID, loggedInUser.ID)
		})
	})
}
