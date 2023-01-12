package http

import (
	"net/http"

	"github.com/bradenrayhorn/beans/beans"
)

type responseAccount struct {
	ID   beans.ID `json:"id"`
	Name string   `json:"name"`
}

func responseFromAccount(account *beans.Account) responseAccount {
	return responseAccount{
		ID:   account.ID,
		Name: string(account.Name),
	}
}

func (s *Server) handleAccountCreate() http.HandlerFunc {
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

		account, err := s.accountContract.Create(r.Context(), getBudgetAuth(r), req.Name)
		if err != nil {
			Error(w, err)
			return
		}

		jsonResponse(w, struct {
			Data response `json:"data"`
		}{
			Data: response{ID: account.ID},
		}, http.StatusOK)
	}
}

func (s *Server) handleAccountsGet() http.HandlerFunc {
	type response struct {
		Data []responseAccount `json:"data"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		accounts, err := s.accountContract.GetAll(r.Context(), getBudgetAuth(r))

		if err != nil {
			Error(w, err)
			return
		}

		res := response{Data: make([]responseAccount, 0, len(accounts))}
		for _, a := range accounts {
			res.Data = append(res.Data, responseFromAccount(a))
		}

		jsonResponse(w, res, http.StatusOK)
	}
}
