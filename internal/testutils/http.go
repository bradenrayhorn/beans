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
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func HTTP(t testing.TB, f http.HandlerFunc, user *beans.User, budget *beans.Budget, body any, status int) string {
	return HTTPWithOptions(t, f, nil, user, budget, body, status)
}

type HTTPOptions struct {
	URLParams map[string]string
}

func HTTPWithOptions(t testing.TB, f http.HandlerFunc, options *HTTPOptions, user *beans.User, budget *beans.Budget, body any, status int) string {
	var reqBody io.Reader
	switch body.(type) {
	case string:
		reqBody = bytes.NewReader([]byte(body.(string)))
	default:
		reqBody = nil
	}
	req := httptest.NewRequest("", "/", reqBody)
	req = req.WithContext(context.WithValue(req.Context(), "budget", budget))
	req = req.WithContext(context.WithValue(req.Context(), "userID", user.ID))

	if options != nil {
		for k, v := range options.URLParams {
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add(k, v)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		}
	}

	w := httptest.NewRecorder()
	f.ServeHTTP(w, req)
	res := w.Result()
	require.Equal(t, status, res.StatusCode)
	data, err := ioutil.ReadAll(res.Body)
	require.Nil(t, err)
	return string(data)
}
