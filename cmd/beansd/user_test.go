package main_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	ta := StartApplication(t)
	defer ta.Stop(t)

	r, err := ta.PostRequest("api/v1/user/register", map[string]interface{}{"username": "user", "password": "pass"})
	require.Nil(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)
}

func TestCannotRegisterWithNoData(t *testing.T) {
	ta := StartApplication(t)
	defer ta.Stop(t)

	r, err := ta.PostRequest("api/v1/user/register", map[string]interface{}{})
	require.Nil(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, r.StatusCode)
}

func TestCannotRegisterSameUsernameTwice(t *testing.T) {
	ta := StartApplication(t)
	defer ta.Stop(t)

	r, err := ta.PostRequest("api/v1/user/register", map[string]interface{}{"username": "user", "password": "pass"})
	require.Nil(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	r, err = ta.PostRequest("api/v1/user/register", map[string]interface{}{"username": "user", "password": "pass"})
	require.Nil(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, r.StatusCode)
}

func TestCannotLoginWithInvalidUsername(t *testing.T) {
	ta := StartApplication(t)
	defer ta.Stop(t)

	r, err := ta.PostRequest("api/v1/user/login", map[string]interface{}{"username": "user", "password": "pass"})
	require.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, r.StatusCode)
	bytes, _ := io.ReadAll(r.Body)
	assert.JSONEq(t, `{"error":"Invalid username or password","code":"unauthorized"}`, string(bytes))
}
