package http

import (
	"net/http"
	"time"

	"github.com/bradenrayhorn/beans/beans"
)

type userResponse struct {
	UserID   string `json:"id"`
	Username string `json:"username"`
}

func (s *Server) handleUserRegister() http.HandlerFunc {
	type request struct {
		Username beans.Username `json:"username"`
		Password beans.Password `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := decodeRequest(r, &req); err != nil {
			Error(w, err)
			return
		}

		err := s.userContract.Register(r.Context(), req.Username, req.Password)
		if err != nil {
			Error(w, err)
			return
		}
	}
}

func (s *Server) handleUserLogin() http.HandlerFunc {
	type request struct {
		Username beans.Username `json:"username"`
		Password beans.Password `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := decodeRequest(r, &req); err != nil {
			Error(w, err)
			return
		}

		session, err := s.userContract.Login(r.Context(), req.Username, req.Password)
		if err != nil {
			Error(w, err)
			return
		}

		cookie := http.Cookie{
			Name:     "session_id",
			Value:    string(session.ID),
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
			Expires:  time.Now().Add(time.Hour * 24 * 30),
		}

		http.SetCookie(w, &cookie)
	}
}

func (s *Server) handleUserLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := s.userContract.Logout(r.Context(), getAuth(r)); err != nil {
			Error(w, beans.WrapError(err, beans.ErrorInternal))
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    "",
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
			Expires:  time.Now().Add(-1 * time.Minute),
		})
	}
}

func (s *Server) handleUserMe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user, err := s.userContract.GetMe(r.Context(), getAuth(r))
		if err != nil {
			Error(w, beans.WrapError(err, beans.ErrorInternal))
			return
		}

		res := userResponse{UserID: user.ID.String(), Username: string(user.Username)}
		jsonResponse(w, res, http.StatusOK)
	}
}
