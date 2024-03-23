package http

import (
	"net/http"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/http/response"
)

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

		err := s.contracts.User.Register(r.Context(), req.Username, req.Password)
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

		session, err := s.contracts.User.Login(r.Context(), req.Username, req.Password)
		if err != nil {
			Error(w, err)
			return
		}

		res := response.Login{Data: response.SessionID{SessionID: session.ID}}
		jsonResponse(w, res, http.StatusOK)
	}
}

func (s *Server) handleUserLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := s.contracts.User.Logout(r.Context(), getAuth(r)); err != nil {
			Error(w, beans.WrapError(err, beans.ErrorInternal))
			return
		}
	}
}

func (s *Server) handleUserMe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user, err := s.contracts.User.GetMe(r.Context(), getAuth(r))
		if err != nil {
			Error(w, beans.WrapError(err, beans.ErrorInternal))
			return
		}

		res := response.GetMe{ID: user.ID, Username: string(user.Username)}
		jsonResponse(w, res, http.StatusOK)
	}
}
