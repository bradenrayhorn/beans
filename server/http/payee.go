package http

import (
	"net/http"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/http/response"
	"github.com/go-chi/chi/v5"
)

func (s *Server) handlePayeeCreate() http.HandlerFunc {
	type request struct {
		Name beans.Name `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := decodeRequest(r, &req); err != nil {
			Error(w, err)
			return
		}

		id, err := s.contracts.Payee.CreatePayee(r.Context(), getBudgetAuth(r), req.Name)
		if err != nil {
			Error(w, err)
			return
		}

		jsonResponse(w, response.CreatePayeeResponse{
			Data: response.ID{ID: id},
		}, http.StatusOK)
	}
}

func (s *Server) handlePayeeGetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payees, err := s.contracts.Payee.GetAll(r.Context(), getBudgetAuth(r))
		if err != nil {
			Error(w, err)
			return
		}

		res := make([]response.Payee, len(payees))
		for i, payee := range payees {
			res[i] = response.Payee{
				ID:   payee.ID,
				Name: payee.Name,
			}
		}

		jsonResponse(w, response.ListPayeesResponse{
			Data: res,
		}, http.StatusOK)
	}
}

func (s *Server) handlePayeeGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := beans.IDFromString(chi.URLParam(r, "payeeID"))
		if err != nil {
			Error(w, beans.WrapError(err, beans.ErrorNotFound))
			return
		}

		payee, err := s.contracts.Payee.Get(r.Context(), getBudgetAuth(r), id)
		if err != nil {
			Error(w, err)
			return
		}

		jsonResponse(w, response.GetPayeeResponse{
			Data: response.Payee{
				ID:   payee.ID,
				Name: payee.Name,
			},
		}, http.StatusOK)
	}
}
