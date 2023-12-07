package http

import (
	"net/http"

	"github.com/bradenrayhorn/beans/server/beans"
)

type payeeResponse struct {
	ID   beans.ID   `json:"id"`
	Name beans.Name `json:"name"`
}

func (s *Server) handlePayeeCreate() http.HandlerFunc {
	type request struct {
		Name beans.Name `json:"name"`
	}
	type response struct {
		ID beans.ID `json:"id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := decodeRequest(r, &req); err != nil {
			Error(w, err)
			return
		}

		payee, err := s.payeeContract.CreatePayee(r.Context(), getBudgetAuth(r), req.Name)
		if err != nil {
			Error(w, err)
			return
		}

		jsonResponse(w, struct {
			Data response `json:"data"`
		}{
			Data: response{ID: payee.ID},
		}, http.StatusOK)
	}
}

func (s *Server) handlePayeeGetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payees, err := s.payeeContract.GetAll(r.Context(), getBudgetAuth(r))
		if err != nil {
			Error(w, err)
			return
		}

		res := make([]payeeResponse, len(payees))
		for i, payee := range payees {
			res[i] = payeeResponse{
				ID:   payee.ID,
				Name: payee.Name,
			}
		}

		jsonResponse(w, struct {
			Data []payeeResponse `json:"data"`
		}{
			Data: res,
		}, http.StatusOK)
	}
}
