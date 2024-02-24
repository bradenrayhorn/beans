package httpadapter

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/http/response"
	"github.com/bradenrayhorn/beans/server/specification"
	"github.com/stretchr/testify/assert"
)

func (a *httpAdapter) UserRegister(t *testing.T, ctx specification.Context, username beans.Username, password beans.Password) error {
	r := a.Request(t, HTTPRequest{
		Method:  "POST",
		Path:    "/api/v1/user/register",
		Body:    fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password),
		Context: ctx,
	})
	return getErrorFromResponse(t, r.Response)
}

func (a *httpAdapter) UserLogin(t *testing.T, ctx specification.Context, username beans.Username, password beans.Password) (beans.SessionID, error) {
	r := a.Request(t, HTTPRequest{
		Method:  "POST",
		Path:    "/api/v1/user/login",
		Body:    fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password),
		Context: ctx,
	})
	if err := getErrorFromResponse(t, r.Response); err != nil {
		return beans.SessionID(""), err
	}

	for _, v := range r.Response.Cookies() {
		if v.Name == "session_id" {
			return beans.SessionID(v.Value), nil
		}
	}

	return beans.SessionID(""), errors.New("http adapter: could not find session id cookie")
}

func (a *httpAdapter) UserLogout(t *testing.T, ctx specification.Context) error {
	r := a.Request(t, HTTPRequest{
		Method:  "POST",
		Path:    "/api/v1/user/logout",
		Context: ctx,
	})
	if err := getErrorFromResponse(t, r.Response); err != nil {
		return err
	}

	for _, v := range r.Response.Cookies() {
		if v.Name == "session_id" {
			assert.Less(t, v.Expires, time.Now(), "http adapter: logout did not expire cookie")
			return nil
		}
	}

	return errors.New("http adapter: could not find session id cookie on logout")
}

func (a *httpAdapter) UserGetMe(t *testing.T, ctx specification.Context) (beans.UserPublic, error) {
	r := a.Request(t, HTTPRequest{
		Method:  "GET",
		Path:    "/api/v1/user/me",
		Context: ctx,
	})
	resp, err := MustParseResponse[response.GetMeResponse](t, r.Response)
	if err != nil {
		return beans.UserPublic{}, err
	}

	return beans.UserPublic{ID: resp.ID, Username: beans.Username(resp.Username)}, nil
}
