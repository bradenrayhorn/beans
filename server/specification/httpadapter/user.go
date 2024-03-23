package httpadapter

import (
	"fmt"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/http/response"
	"github.com/bradenrayhorn/beans/server/specification"
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
	resp, err := MustParseResponse[response.Login](t, r.Response)
	if err != nil {
		return beans.SessionID(""), err
	}

	return resp.Data.SessionID, nil
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

	return nil
}

func (a *httpAdapter) UserGetMe(t *testing.T, ctx specification.Context) (beans.UserPublic, error) {
	r := a.Request(t, HTTPRequest{
		Method:  "GET",
		Path:    "/api/v1/user/me",
		Context: ctx,
	})
	resp, err := MustParseResponse[response.GetMe](t, r.Response)
	if err != nil {
		return beans.UserPublic{}, err
	}

	return beans.UserPublic{ID: resp.ID, Username: beans.Username(resp.Username)}, nil
}
