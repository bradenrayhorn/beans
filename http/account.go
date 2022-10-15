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

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := decodeRequest(r, &req); err != nil {
			Error(w, err)
			return
		}

		_, err := s.accountService.Create(r.Context(), req.Name, getBudget(r).ID)
		if err != nil {
			Error(w, err)
			return
		}
	}
}

func (s *Server) handleAccountsGet() http.HandlerFunc {
	type response struct {
		Data []responseAccount `json:"data"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		accounts, err := s.accountRepository.GetForBudget(r.Context(), getBudget(r).ID)

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
