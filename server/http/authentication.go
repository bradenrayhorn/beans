package http

import (
	"context"
	"net/http"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/http/httpcontext"
)

func (s *Server) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Header.Get("Authorization")

		authCtx, err := s.services.User.GetAuth(r.Context(), beans.SessionID(sessionID))
		if err != nil {
			Error(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), httpcontext.Auth, authCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getAuth(r *http.Request) *beans.AuthContext {
	return r.Context().Value(httpcontext.Auth).(*beans.AuthContext)
}
