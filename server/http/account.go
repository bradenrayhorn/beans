package http

import (
	"net/http"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/http/response"
)

func (s *Server) handleAccountCreate() http.HandlerFunc {
	type request struct {
		Name beans.Name `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := decodeRequest(r, &req); err != nil {
			Error(w, err)
			return
		}

		account, err := s.accountContract.Create(r.Context(), getBudgetAuth(r), req.Name)
		if err != nil {
			Error(w, err)
			return
		}

		jsonResponse(w, response.CreateAccountResponse{
			Data: response.ID{ID: account.ID},
		}, http.StatusOK)
	}
}

func (s *Server) handleAccountsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accounts, err := s.accountContract.GetAll(r.Context(), getBudgetAuth(r))

		if err != nil {
			Error(w, err)
			return
		}

		res := make([]response.ListAccount, 0, len(accounts))
		for _, a := range accounts {
			res = append(res, response.ListAccount{ID: a.ID, Name: string(a.Name), Balance: a.Balance})
		}

		jsonResponse(w, response.ListAccountResponse{Data: res}, http.StatusOK)
	}
}
