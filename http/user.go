package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bradenrayhorn/beans/beans"
)

func (s *Server) handleUserRegister() http.HandlerFunc {
	type request struct {
		Username beans.Username `json:"username"`
		Password beans.Password `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		_, err := s.userService.CreateUser(r.Context(), req.Username, req.Password)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) handleUserLogin() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {}
}
