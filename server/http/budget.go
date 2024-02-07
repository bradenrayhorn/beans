package http

import (
	"context"
	"net/http"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/http/httpcontext"
	"github.com/bradenrayhorn/beans/server/http/response"
	"github.com/go-chi/chi/v5"
)

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

		budget, err := s.contracts.Budget.Create(r.Context(), getAuth(r), req.Name)
		if err != nil {
			Error(w, err)
			return
		}

		jsonResponse(w, response.CreateBudgetResponse{
			Data: response.ID{ID: budget.ID}},
			http.StatusOK)
	}
}

func (s *Server) handleBudgetGetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		budgets, err := s.contracts.Budget.GetAll(r.Context(), getAuth(r))

		if err != nil {
			Error(w, err)
			return
		}

		res := response.ListBudgetsResposne{Data: []response.Budget{}}
		for _, b := range budgets {
			res.Data = append(res.Data, response.Budget{ID: b.ID, Name: b.Name})
		}

		jsonResponse(w, res, http.StatusOK)
	}
}

func (s *Server) handleBudgetGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		budgetID, err := beans.BeansIDFromString(chi.URLParam(r, "budgetID"))
		if err != nil {
			Error(w, beans.WrapError(err, beans.ErrorNotFound))
			return
		}

		budget, err := s.contracts.Budget.Get(r.Context(), getAuth(r), budgetID)
		if err != nil {
			Error(w, err)
			return
		}

		res := response.GetBudgetResponse{Data: response.Budget{
			ID: budget.ID, Name: budget.Name,
		}}

		jsonResponse(w, res, http.StatusOK)
	}
}

// middleware

func (s *Server) parseBudgetHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		budgetID, err := beans.BeansIDFromString(r.Header.Get("Budget-ID"))
		if err != nil {
			Error(w, beans.WrapError(err, beans.ErrorNotFound))
			return
		}

		budget, err := s.contracts.Budget.Get(r.Context(), getAuth(r), budgetID)
		if err != nil {
			Error(w, err)
			return
		}

		auth, err := beans.NewBudgetAuthContext(getAuth(r), budget)
		if err != nil {
			Error(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), httpcontext.Budget, budget)
		ctx = context.WithValue(ctx, httpcontext.BudgetAuth, auth)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getBudgetAuth(r *http.Request) *beans.BudgetAuthContext {
	return r.Context().Value(httpcontext.BudgetAuth).(*beans.BudgetAuthContext)
}
