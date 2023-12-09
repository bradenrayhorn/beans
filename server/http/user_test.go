package http_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser(t *testing.T) {
	user := &beans.User{ID: beans.NewBeansID(), Username: "user"}

	t.Run("register", func(t *testing.T) {
		test := newHttpTest(t)
		defer test.Stop(t)

		test.userContract.RegisterFunc.PushReturn(nil)

		res := test.DoRequest(t, HTTPRequest{
			method: "POST",
			path:   "/api/v1/user/register",
			body:   `{"username":"user","password":"password"}`,
		})

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Empty(t, res.body)

		params := test.userContract.RegisterFunc.History()[0]
		assert.Equal(t, "user", string(params.Arg1))
		assert.Equal(t, "password", string(params.Arg2))
	})

	t.Run("login", func(t *testing.T) {
		test := newHttpTest(t)
		defer test.Stop(t)

		test.userContract.LoginFunc.PushReturn(&beans.Session{
			ID: beans.SessionID("12345"),
		}, nil)

		res := test.DoRequest(t, HTTPRequest{
			method: "POST",
			path:   "/api/v1/user/login",
			body:   `{"username":"user","password":"password"}`,
		})

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Empty(t, res.body)
		var cookie *http.Cookie
		for _, c := range res.Cookies() {
			if c.Name == "session_id" {
				cookie = c
			}
		}
		require.NotNil(t, cookie)
		assert.Equal(t, "12345", cookie.Value)

		params := test.userContract.LoginFunc.History()[0]
		assert.Equal(t, "user", string(params.Arg1))
		assert.Equal(t, "password", string(params.Arg2))
	})

	t.Run("logout", func(t *testing.T) {
		test := newHttpTest(t)
		defer test.Stop(t)

		test.userContract.LogoutFunc.PushReturn(nil)

		res := test.DoRequest(t, HTTPRequest{
			method: "POST",
			path:   "/api/v1/user/logout",
			user:   user,
		})

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Empty(t, res.body)
		var cookie *http.Cookie
		for _, c := range res.Cookies() {
			if c.Name == "session_id" {
				cookie = c
			}
		}
		require.NotNil(t, cookie)
		assert.LessOrEqual(t, cookie.Expires, time.Now())

		params := test.userContract.LogoutFunc.History()[0]
		assert.Equal(t, user.ID, params.Arg1.UserID())
	})

	t.Run("get me", func(t *testing.T) {
		test := newHttpTest(t)
		defer test.Stop(t)

		test.userContract.GetMeFunc.PushReturn(user, nil)

		res := test.DoRequest(t, HTTPRequest{
			method: "GET",
			path:   "/api/v1/user/me",
			user:   user,
		})

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.JSONEq(t, fmt.Sprintf(
			`{"username":"user","id":"%s"}`,
			user.ID.String(),
		), res.body)

		params := test.userContract.GetMeFunc.History()[0]
		assert.Equal(t, user.ID, params.Arg1.UserID())
	})
}
