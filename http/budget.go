package http

import (
	"context"
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
		Name beans.Name `json:"name"`
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
		budget := s.getBudget(chi.URLParam(r, "budgetID"), w, r)
		if budget == nil {
			return
		}

		res := response{Data: responseBudget{ID: budget.ID.String(), Name: string(budget.Name)}}

		jsonResponse(w, res, http.StatusOK)
	}
}

// middleware

func (s *Server) parseBudgetHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		budget := s.getBudget(r.Header.Get("Budget-ID"), w, r)
		if budget == nil {
			return
		}

		ctx := context.WithValue(r.Context(), "budget", budget)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) getBudget(id string, w http.ResponseWriter, r *http.Request) *beans.Budget {
	budgetID, err := beans.BeansIDFromString(id)
	if err != nil {
		Error(w, beans.WrapError(err, beans.ErrorNotFound))
		return nil
	}

	budget, err := s.budgetRepository.Get(r.Context(), budgetID)

	if err != nil {
		Error(w, err)
		return nil
	}

	if !budget.UserHasAccess(getUserID(r)) {
		Error(w, beans.ErrorForbidden)
		return nil
	}

	return budget
}

func getBudget(r *http.Request) *beans.Budget {
	return r.Context().Value("budget").(*beans.Budget)
}
