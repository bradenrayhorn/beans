package http

import (
	"net/http"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/http/request"
	"github.com/bradenrayhorn/beans/server/http/response"
	"github.com/go-chi/chi/v5"
)

func (s *Server) handleAccountCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request.CreateAccount
		if err := decodeRequest(r, &req); err != nil {
			Error(w, err)
			return
		}

		accountID, err := s.contracts.Account.Create(r.Context(), getBudgetAuth(r), beans.AccountCreate{
			Name:      req.Name,
			OffBudget: req.OffBudget,
		})
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
			res = append(res, response.ListAccount{
				ID:        a.ID,
				Name:      string(a.Name),
				Balance:   a.Balance,
				OffBudget: a.OffBudget,
			})
		}

		jsonResponse(w, response.ListAccountResponse{Data: res}, http.StatusOK)
	}
}

func (s *Server) handleAccountsGetTransactable() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accounts, err := s.contracts.Account.GetTransactable(r.Context(), getBudgetAuth(r))

		if err != nil {
			Error(w, err)
			return
		}

		res := make([]response.Account, 0, len(accounts))
		for _, a := range accounts {
			res = append(res, response.Account{ID: a.ID, Name: string(a.Name), OffBudget: a.OffBudget})
		}

		jsonResponse(w, response.GetTransactableAccounts{Data: res}, http.StatusOK)
	}
}

func (s *Server) handleAccountGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID, err := beans.IDFromString(chi.URLParam(r, "accountID"))
		if err != nil {
			Error(w, beans.WrapError(err, beans.ErrorNotFound))
			return
		}

		account, err := s.contracts.Account.Get(r.Context(), getBudgetAuth(r), accountID)

		if err != nil {
			Error(w, err)
			return
		}

		jsonResponse(w, response.GetAccountResponse{Data: response.Account{
			ID:        account.ID,
			Name:      string(account.Name),
			OffBudget: account.OffBudget,
		}}, http.StatusOK)
	}
}
