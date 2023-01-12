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
	type response struct {
		Data responseBudget `json:"data"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := decodeRequest(r, &req); err != nil {
			Error(w, err)
			return
		}

		budget, err := s.budgetContract.Create(r.Context(), getAuth(r), req.Name)
		if err != nil {
			Error(w, err)
			return
		}

		jsonResponse(w, response{Data: responseBudget{ID: budget.ID.String(), Name: string(budget.Name)}}, http.StatusOK)
	}
}

func (s *Server) handleBudgetGetAll() http.HandlerFunc {
	type response struct {
		Data []responseBudget `json:"data"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		budgets, err := s.budgetContract.GetAll(r.Context(), getAuth(r))

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
	type responseData struct {
		LatestMonth beans.ID `json:"latest_month_id"`
		responseBudget
	}
	type response struct {
		Data responseData `json:"data"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		budgetID, err := beans.BeansIDFromString(chi.URLParam(r, "budgetID"))
		if err != nil {
			Error(w, beans.WrapError(err, beans.ErrorNotFound))
			return
		}

		budget, latestMonth, err := s.budgetContract.Get(r.Context(), getAuth(r), budgetID)
		if err != nil {
			Error(w, err)
			return
		}

		res := response{Data: responseData{
			responseBudget: responseBudget{ID: budget.ID.String(), Name: string(budget.Name)},
			LatestMonth:    latestMonth.ID,
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

		budget, err := s.budgetRepository.Get(r.Context(), budgetID)
		if err != nil {
			Error(w, err)
			return
		}

		auth, err := beans.NewBudgetAuthContext(getAuth(r), budget)
		if err != nil {
			Error(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), "budget", budget)
		ctx = context.WithValue(ctx, "budget_auth", auth)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getBudget(r *http.Request) *beans.Budget {
	return r.Context().Value("budget").(*beans.Budget)
}

func getBudgetAuth(r *http.Request) *beans.BudgetAuthContext {
	return r.Context().Value("budget_auth").(*beans.BudgetAuthContext)
}
