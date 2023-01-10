package http

import (
	"context"
	"net/http"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/go-chi/chi/v5"
)

func (s *Server) handleMonthGet() http.HandlerFunc {
	type responseCategory struct {
		ID         beans.ID     `json:"id"`
		Assigned   beans.Amount `json:"assigned"`
		Spent      beans.Amount `json:"spent"`
		CategoryID beans.ID     `json:"category_id"`
	}
	type responseMonth struct {
		ID         beans.ID           `json:"id"`
		Date       string             `json:"date"`
		Categories []responseCategory `json:"categories"`
	}
	type response struct {
		Data responseMonth `json:"data"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		monthID := getMonth(r).ID

		month, categories, err := s.monthContract.Get(r.Context(), monthID)
		if err != nil {
			Error(w, err)
			return
		}

		responseCategories := make([]responseCategory, len(categories))
		for i, category := range categories {
			responseCategories[i] = responseCategory{
				ID:         category.ID,
				Assigned:   category.Amount,
				Spent:      category.Spent,
				CategoryID: category.CategoryID,
			}
		}

		jsonResponse(w, response{
			Data: responseMonth{
				ID:         month.ID,
				Date:       month.Date.String(),
				Categories: responseCategories,
			},
		}, http.StatusOK)
	}
}

func (s *Server) handleMonthCreate() http.HandlerFunc {
	type request struct {
		Date beans.Date `json:"date"`
	}

	type response struct {
		ID beans.ID `json:"month_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := decodeRequest(r, &req); err != nil {
			Error(w, err)
			return
		}

		month, err := s.monthContract.CreateMonth(r.Context(), getBudget(r).ID, beans.NewMonthDate(req.Date))
		if err != nil {
			Error(w, err)
			return
		}

		jsonResponse(w, struct {
			Data response `json:"data"`
		}{
			Data: response{ID: month.ID},
		}, http.StatusOK)
	}
}

func (s *Server) handleMonthCategoryUpdate() http.HandlerFunc {
	type request struct {
		CategoryID beans.ID     `json:"category_id"`
		Amount     beans.Amount `json:"amount"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := decodeRequest(r, &req); err != nil {
			Error(w, err)
			return
		}

		monthID := getMonth(r).ID

		if err := s.monthContract.SetCategoryAmount(r.Context(), monthID, req.CategoryID, req.Amount); err != nil {
			Error(w, err)
			return
		}
	}
}

func (s *Server) validateMonth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		monthID, err := beans.BeansIDFromString(chi.URLParam(r, "monthID"))
		if err != nil {
			Error(w, err)
			return
		}

		month, err := s.monthRepository.Get(r.Context(), monthID)
		if err != nil {
			Error(w, err)
			return
		}

		// month must be a part of the budget
		if month.BudgetID != getBudget(r).ID {
			Error(w, beans.ErrorNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "month", month)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getMonth(r *http.Request) *beans.Month {
	return r.Context().Value("month").(*beans.Month)
}
