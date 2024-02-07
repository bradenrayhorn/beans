package httpadapter

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/specification"
)

func (a *HTTPAdapter) UserRegister(t *testing.T, ctx specification.Context, username beans.Username, password beans.Password) error {
	r := a.Request(t, HTTPRequest{
		Method:  "POST",
		Path:    "/api/v1/user/register",
		Body:    fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password),
		Context: ctx,
	})
	return getErrorFromResponse(t, r.Response)
}

func (a *HTTPAdapter) UserLogin(t *testing.T, ctx specification.Context, username beans.Username, password beans.Password) (beans.SessionID, error) {
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
