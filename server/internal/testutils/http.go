package testutils

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/http/httpcontext"
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
	switch body := body.(type) {
	case string:
		reqBody = bytes.NewReader([]byte(body))
	default:
		reqBody = nil
	}
	req := httptest.NewRequest("", "/", reqBody)
	req = req.WithContext(context.WithValue(req.Context(), httpcontext.Budget, budget))

	auth := beans.NewAuthContext(user.ID, beans.SessionID("1234"))
	req = req.WithContext(context.WithValue(req.Context(), httpcontext.Auth, auth))

	if budget != nil {
		budgetAuth, err := beans.NewBudgetAuthContext(auth, budget)
		require.Nil(t, err)
		req = req.WithContext(context.WithValue(req.Context(), httpcontext.BudgetAuth, budgetAuth))
	}

	if options != nil {
		rctx := chi.NewRouteContext()
		for k, v := range options.URLParams {
			rctx.URLParams.Add(k, v)
		}
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	}

	w := httptest.NewRecorder()
	f.ServeHTTP(w, req)
	res := w.Result()
	require.Equal(t, status, res.StatusCode)
	data, err := io.ReadAll(res.Body)
	require.Nil(t, err)
	return string(data)
}
