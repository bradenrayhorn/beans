package http

import (
	"net/http"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/go-chi/chi/v5"
)

func (s *Server) handleMonthGet() http.HandlerFunc {
	type responseCategory struct {
		ID         beans.ID     `json:"id"`
		Assigned   beans.Amount `json:"assigned"`
		Activity   beans.Amount `json:"activity"`
		CategoryID beans.ID     `json:"category_id"`
	}
	type responseMonth struct {
		ID         beans.ID           `json:"id"`
		Date       string             `json:"date"`
		Budgetable beans.Amount       `json:"budgetable"`
		Categories []responseCategory `json:"categories"`
	}
	type response struct {
		Data responseMonth `json:"data"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		monthID, err := beans.BeansIDFromString(chi.URLParam(r, "monthID"))
		if err != nil {
			Error(w, err)
			return
		}

		month, categories, budgetable, err := s.monthContract.Get(r.Context(), getBudgetAuth(r), monthID)
		if err != nil {
			Error(w, err)
			return
		}

		responseCategories := make([]responseCategory, len(categories))
		for i, category := range categories {
			responseCategories[i] = responseCategory{
				ID:         category.ID,
				Assigned:   category.Amount,
				Activity:   category.Activity,
				CategoryID: category.CategoryID,
			}
		}

		jsonResponse(w, response{
			Data: responseMonth{
				ID:         month.ID,
				Date:       month.Date.String(),
				Budgetable: budgetable,
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

		month, err := s.monthContract.CreateMonth(r.Context(), getBudgetAuth(r), beans.NewMonthDate(req.Date))
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

		monthID, err := beans.BeansIDFromString(chi.URLParam(r, "monthID"))
		if err != nil {
			Error(w, err)
			return
		}

		if err := s.monthContract.SetCategoryAmount(r.Context(), getBudgetAuth(r), monthID, req.CategoryID, req.Amount); err != nil {
			Error(w, err)
			return
		}
	}
}
