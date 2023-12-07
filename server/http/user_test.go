package http

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/mocks"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	contract := mocks.NewMockUserContract()
	sv := Server{userContract: contract}

	user := &beans.User{ID: beans.NewBeansID(), Username: "user"}

	t.Run("register", func(t *testing.T) {
		contract.RegisterFunc.PushReturn(nil)

		req := `{"username":"user","password":"password"}`
		res := testutils.HTTP(t, sv.handleUserRegister(), user, nil, req, http.StatusOK)

		assert.Equal(t, "", res)

		params := contract.RegisterFunc.History()[0]
		assert.Equal(t, "user", string(params.Arg1))
		assert.Equal(t, "password", string(params.Arg2))
	})

	t.Run("login", func(t *testing.T) {
		contract.LoginFunc.PushReturn(&beans.Session{}, nil)

		req := `{"username":"user","password":"password"}`
		res := testutils.HTTP(t, sv.handleUserLogin(), user, nil, req, http.StatusOK)

		assert.Equal(t, "", res)

		params := contract.LoginFunc.History()[0]
		assert.Equal(t, "user", string(params.Arg1))
		assert.Equal(t, "password", string(params.Arg2))
	})

	t.Run("logout", func(t *testing.T) {
		contract.LogoutFunc.PushReturn(nil)

		res := testutils.HTTP(t, sv.handleUserLogout(), user, nil, nil, http.StatusOK)

		assert.Equal(t, "", res)

		params := contract.LogoutFunc.History()[0]
		assert.Equal(t, user.ID, params.Arg1.UserID())
	})

	t.Run("get me", func(t *testing.T) {
		contract.GetMeFunc.PushReturn(user, nil)

		res := testutils.HTTP(t, sv.handleUserMe(), user, nil, nil, http.StatusOK)

		assert.JSONEq(t, fmt.Sprintf(`{"username":"user","id":"%s"}`, user.ID.String()), res)

		params := contract.GetMeFunc.History()[0]
		assert.Equal(t, user.ID, params.Arg1.UserID())
	})
}
