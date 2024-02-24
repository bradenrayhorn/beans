package specification

import (
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testUser(t *testing.T, interactor Interactor) {

	t.Run("register", func(t *testing.T) {

		t.Run("username is required", func(t *testing.T) {
			err := interactor.UserRegister(t, Context{}, beans.Username(""), beans.Password("no"))
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("password is required", func(t *testing.T) {
			err := interactor.UserRegister(t, Context{}, beans.Username("user"), beans.Password(""))
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("cannot register with taken username", func(t *testing.T) {
			username := beans.Username(beans.NewID().String())
			err := interactor.UserRegister(t, Context{}, username, beans.Password("pass"))
			require.NoError(t, err)

			// register a second time
			err = interactor.UserRegister(t, Context{}, username, beans.Password("pass"))
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})
	})

	t.Run("register, login, and get me", func(t *testing.T) {
		username := beans.Username(beans.NewID().String())

		// register
		err := interactor.UserRegister(t, Context{}, username, beans.Password("pass"))
		require.NoError(t, err)

		// login
		sessionID, err := interactor.UserLogin(t, Context{}, username, beans.Password("pass"))
		require.NoError(t, err)

		// get me and verify results
		me, err := interactor.UserGetMe(t, Context{SessionID: sessionID})
		require.NoError(t, err)

		assert.Equal(t, false, me.ID.Empty())
		assert.Equal(t, username, me.Username)
	})

	t.Run("login", func(t *testing.T) {

		t.Run("username is required", func(t *testing.T) {
			_, err := interactor.UserLogin(t, Context{}, beans.Username(""), beans.Password("no"))
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("password is required", func(t *testing.T) {
			_, err := interactor.UserLogin(t, Context{}, beans.Username("user"), beans.Password(""))
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("cannot login with non existent username", func(t *testing.T) {
			username := beans.Username(beans.NewID().String())
			_, err := interactor.UserLogin(t, Context{}, username, beans.Password("pass"))
			testutils.AssertErrorCode(t, err, beans.EUNAUTHORIZED)
		})

		t.Run("cannot login with invalid password", func(t *testing.T) {
			username := beans.Username(beans.NewID().String())

			// register user
			err := interactor.UserRegister(t, Context{}, username, beans.Password("pass"))
			require.NoError(t, err)

			// try to login with other password
			_, err = interactor.UserLogin(t, Context{}, username, beans.Password("pass-bad"))
			testutils.AssertErrorCode(t, err, beans.EUNAUTHORIZED)

			// login with good password
			_, err = interactor.UserLogin(t, Context{}, username, beans.Password("pass"))
			require.NoError(t, err)
		})
	})

	t.Run("logout", func(t *testing.T) {
		username := beans.Username(beans.NewID().String())

		// register user
		err := interactor.UserRegister(t, Context{}, username, beans.Password("pass"))
		require.NoError(t, err)

		// login
		sessionID, err := interactor.UserLogin(t, Context{}, username, beans.Password("pass"))
		require.NoError(t, err)
		ctx := Context{SessionID: sessionID}

		// should be able to get me
		_, err = interactor.UserGetMe(t, ctx)
		require.NoError(t, err)

		// logout
		err = interactor.UserLogout(t, ctx)
		require.NoError(t, err)

		// get me should fail now
		_, err = interactor.UserGetMe(t, ctx)
		testutils.AssertErrorCode(t, err, beans.EUNAUTHORIZED)
	})
}
