package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/http/response"
	"github.com/go-chi/chi/v5"
)

func (s *Server) handleMonthGetOrCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dateParam := chi.URLParam(r, "date")
		var date beans.Date
		if err := json.Unmarshal([]byte(fmt.Sprintf(`"%s"`, dateParam)), &date); err != nil {
			Error(w, beans.WrapError(err, beans.ErrorInvalid))
			return
		}

		month, categories, budgetable, err := s.contracts.Month.GetOrCreate(r.Context(), getBudgetAuth(r), beans.NewMonthDate(date))
		if err != nil {
			Error(w, err)
			return
		}

		responseCategories := make([]response.MonthCategory, len(categories))
		for i, category := range categories {
			responseCategories[i] = response.MonthCategory{
				ID:         category.ID,
				Assigned:   category.Amount,
				Activity:   category.Activity,
				Available:  category.Available,
				CategoryID: category.CategoryID,
			}
		}

		jsonResponse(w, response.GetMonthResponse{
			Data: response.Month{
				ID:          month.ID,
				Date:        month.Date.String(),
				Budgetable:  budgetable,
				Carryover:   month.Carryover,
				Income:      month.Income,
				Assigned:    month.Assigned,
				CarriedOver: month.CarriedOver,
				Categories:  responseCategories,
			},
		}, http.StatusOK)
	}
}

func (s *Server) handleMonthUpdate() http.HandlerFunc {
	type request struct {
		Amount beans.Amount `json:"carryover"`
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

		if err := s.contracts.Month.Update(r.Context(), getBudgetAuth(r), monthID, req.Amount); err != nil {
			Error(w, err)
			return
		}
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

		if err := s.contracts.Month.SetCategoryAmount(r.Context(), getBudgetAuth(r), monthID, req.CategoryID, req.Amount); err != nil {
			Error(w, err)
			return
		}
	}
}
