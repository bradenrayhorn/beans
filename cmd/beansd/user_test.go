package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCanRegisterAndLoginAndGetMe(t *testing.T) {
	ta := StartApplication(t)
	defer ta.Stop(t)

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
}

func TestCannotRegisterWithNoData(t *testing.T) {
	ta := StartApplication(t)
	defer ta.Stop(t)

	r := ta.PostRequest(t, "api/v1/user/register", nil)
	assert.Equal(t, http.StatusUnprocessableEntity, r.StatusCode)
}

func TestCannotRegisterSameUsernameTwice(t *testing.T) {
	ta := StartApplication(t)
	defer ta.Stop(t)

	ta.CreateTestUser(t, "user", "password")

	r := ta.PostRequest(t, "api/v1/user/register", &RequestOptions{Body: `{"username": "user", "password": "password"}`})
	assert.Equal(t, http.StatusUnprocessableEntity, r.StatusCode)
}

func TestCanLogin(t *testing.T) {
	ta := StartApplication(t)
	defer ta.Stop(t)

	ta.CreateTestUser(t, "user", "password")
	r := ta.PostRequest(t, "api/v1/user/login", &RequestOptions{Body: `{"username": "user", "password": "password"}`})
	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.NotEmpty(t, r.SessionIDFromCookie)
}

func TestCannotLoginWithInvalidUsername(t *testing.T) {
	ta := StartApplication(t)
	defer ta.Stop(t)

	r := ta.PostRequest(t, "api/v1/user/login", &RequestOptions{Body: `{"username": "user", "password": "password"}`})
	assert.Equal(t, http.StatusUnauthorized, r.StatusCode)
	assert.JSONEq(t, `{"error":"Invalid username or password","code":"unauthorized"}`, r.Body)
}

func TestCannotGetMeWithNoSession(t *testing.T) {
	ta := StartApplication(t)
	defer ta.Stop(t)

	r := ta.GetRequest(t, "api/v1/user/me", nil)
	assert.Equal(t, http.StatusUnauthorized, r.StatusCode)
}
