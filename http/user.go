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

		_, err := s.userContract.CreateUser(r.Context(), req.Username, req.Password)
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

		user, err := s.userContract.Login(r.Context(), req.Username, req.Password)
		if err != nil {
			Error(w, err)
			return
		}

		session, err := s.sessionRepository.Create(user.ID)
		if err != nil {
			Error(w, beans.WrapError(err, beans.ErrorInternal))
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

		res := userResponse{UserID: user.ID.String(), Username: string(user.Username)}
		jsonResponse(w, res, http.StatusOK)
	}
}

func (s *Server) handleUserLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			Error(w, beans.WrapError(err, beans.ErrorInternal))
			return
		}

		err = s.sessionRepository.Delete(beans.SessionID(cookie.Value))
		if err != nil {
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
		userID := getUserID(r)

		user, err := s.userRepository.Get(r.Context(), userID)
		if err != nil {
			Error(w, beans.WrapError(err, beans.ErrorInternal))
			return
		}

		res := userResponse{UserID: userID.String(), Username: string(user.Username)}
		jsonResponse(w, res, http.StatusOK)
	}
}
