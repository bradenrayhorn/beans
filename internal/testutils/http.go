package testutils

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/stretchr/testify/require"
)

func HTTP(t testing.TB, f http.HandlerFunc, user *beans.User, budget *beans.Budget, body any, status int) string {
	return HTTPWithContext(t, f, func(ctx context.Context) context.Context { return ctx }, user, budget, body, status)
}

func HTTPWithContext(t testing.TB, f http.HandlerFunc, withContext func(context.Context) context.Context, user *beans.User, budget *beans.Budget, body any, status int) string {
	var reqBody io.Reader
	switch body.(type) {
	case string:
		reqBody = bytes.NewReader([]byte(body.(string)))
	default:
		reqBody = nil
	}
	req := httptest.NewRequest("", "/", reqBody)
	req = req.WithContext(withContext(req.Context()))
	req = req.WithContext(context.WithValue(req.Context(), "budget", budget))
	req = req.WithContext(context.WithValue(req.Context(), "userID", user.ID))
	w := httptest.NewRecorder()
	f.ServeHTTP(w, req)
	res := w.Result()
	require.Equal(t, status, res.StatusCode)
	data, err := ioutil.ReadAll(res.Body)
	require.Nil(t, err)
	return string(data)
}
