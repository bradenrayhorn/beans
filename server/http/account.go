package http

import (
	"net/http"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/http/response"
	"github.com/go-chi/chi/v5"
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

		accountID, err := s.contracts.Account.Create(r.Context(), getBudgetAuth(r), req.Name)
		if err != nil {
			Error(w, err)
			return
		}

		jsonResponse(w, response.CreateAccountResponse{
			Data: response.ID{ID: accountID},
		}, http.StatusOK)
	}
}

func (s *Server) handleAccountsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accounts, err := s.contracts.Account.GetAll(r.Context(), getBudgetAuth(r))

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

func (s *Server) handleAccountGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID, err := beans.BeansIDFromString(chi.URLParam(r, "accountID"))
		if err != nil {
			Error(w, beans.WrapError(err, beans.ErrorNotFound))
			return
		}

		account, err := s.contracts.Account.Get(r.Context(), getBudgetAuth(r), accountID)

		if err != nil {
			Error(w, err)
			return
		}

		jsonResponse(w, response.GetAccountResponse{Data: response.ListAccount{
			ID:      account.ID,
			Name:    string(account.Name),
			Balance: account.Balance,
		}}, http.StatusOK)
	}
}
