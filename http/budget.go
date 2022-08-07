package http

import (
	"net/http"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/go-chi/chi/v5"
)

type responseBudget struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (s *Server) handleBudgetCreate() http.HandlerFunc {
	type request struct {
		Name beans.BudgetName `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := decodeRequest(r, &req); err != nil {
			Error(w, err)
			return
		}

		_, err := s.budgetService.CreateBudget(r.Context(), req.Name, getUserID(r))
		if err != nil {
			Error(w, err)
			return
		}
	}
}

func (s *Server) handleBudgetGetAll() http.HandlerFunc {
	type response struct {
		Data []responseBudget `json:"data"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		budgets, err := s.budgetRepository.GetBudgetsForUser(r.Context(), getUserID(r))

		if err != nil {
			Error(w, err)
			return
		}

		res := response{Data: []responseBudget{}}
		for _, b := range budgets {
			res.Data = append(res.Data, responseBudget{ID: b.ID.String(), Name: string(b.Name)})
		}

		jsonResponse(w, res, http.StatusOK)
	}
}

func (s *Server) handleBudgetGet() http.HandlerFunc {
	type response struct {
		Data responseBudget `json:"data"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		budgetID, err := beans.BeansIDFromString(chi.URLParam(r, "budgetID"))
		if err != nil {
			Error(w, beans.WrapError(err, beans.ErrorNotFound))
			return
		}

		budget, err := s.budgetRepository.Get(r.Context(), budgetID)

		if err != nil {
			Error(w, err)
			return
		}

		if !budget.UserHasAccess(getUserID(r)) {
			Error(w, beans.ErrorForbidden)
			return
		}

		res := response{Data: responseBudget{ID: budgetID.String(), Name: string(budget.Name)}}

		jsonResponse(w, res, http.StatusOK)
	}
}
