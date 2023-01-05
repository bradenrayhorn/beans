package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsers(t *testing.T) {
	ta := StartApplication(t)
	defer ta.Stop(t)

	t.Run("can register, login, get me, and logout", func(t *testing.T) {
		r := ta.PostRequest(t, "api/v1/user/register", &RequestOptions{Body: `{"username": "user", "password": "password"}`})
		assert.Equal(t, http.StatusOK, r.StatusCode)

		r = ta.PostRequest(t, "api/v1/user/login", &RequestOptions{Body: `{"username": "user", "password": "password"}`})
		assert.Equal(t, http.StatusOK, r.StatusCode)
		sessionID := r.SessionIDFromCookie

		r = ta.GetRequest(t, "api/v1/user/me", &RequestOptions{SessionID: sessionID})
		assert.Equal(t, http.StatusOK, r.StatusCode)
		type response struct {
			Username string `json:"username"`
		}
		var responseJson response
		err := json.NewDecoder(bytes.NewReader([]byte(r.Body))).Decode(&responseJson)
		require.Nil(t, err)
		assert.Equal(t, "user", responseJson.Username)

		r = ta.PostRequest(t, "api/v1/user/logout", &RequestOptions{SessionID: sessionID})
		assert.Equal(t, http.StatusOK, r.StatusCode)

		var sessionCookie *http.Cookie
		for _, c := range r.resp.Cookies() {
			if c.Name == "session_id" {
				sessionCookie = c
			}
		}

		require.NotNil(t, sessionCookie)
		assert.Less(t, sessionCookie.Expires, time.Now())
		assert.Equal(t, "", sessionCookie.Value)
	})

	t.Run("cannot register with no data", func(t *testing.T) {
		r := ta.PostRequest(t, "api/v1/user/register", nil)
		assert.Equal(t, http.StatusUnprocessableEntity, r.StatusCode)
	})

	t.Run("cannot register same username twice", func(t *testing.T) {
		ta.CreateUser(t, "user3", "password")

		r := ta.PostRequest(t, "api/v1/user/register", &RequestOptions{Body: `{"username": "user3", "password": "password"}`})
		assert.Equal(t, http.StatusUnprocessableEntity, r.StatusCode)
	})

	t.Run("cannot login with invalid username", func(t *testing.T) {
		r := ta.PostRequest(t, "api/v1/user/login", &RequestOptions{Body: `{"username": "user4", "password": "password"}`})
		assert.Equal(t, http.StatusUnauthorized, r.StatusCode)
		assert.JSONEq(t, `{"error":"Invalid username or password","code":"unauthorized"}`, r.Body)
	})

	t.Run("cannot get me with no session", func(t *testing.T) {
		r := ta.GetRequest(t, "api/v1/user/me", nil)
		assert.Equal(t, http.StatusUnauthorized, r.StatusCode)
	})
}
