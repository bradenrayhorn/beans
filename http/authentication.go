package http

import (
	"context"
	"net/http"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/http/httpcontext"
)

func (s *Server) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			Error(w, beans.ErrorUnauthorized)
			return
		}

		session, err := s.sessionRepository.Get(beans.SessionID(cookie.Value))
		if err != nil {
			Error(w, beans.WrapError(err, beans.ErrorUnauthorized))
			return
		}

		ctx := context.WithValue(r.Context(), httpcontext.Auth, beans.NewAuthContext(session.UserID, session.ID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getAuth(r *http.Request) *beans.AuthContext {
	return r.Context().Value(httpcontext.Auth).(*beans.AuthContext)
}
