package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/http/request"
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

		month, err := s.contracts.Month.GetOrCreate(r.Context(), getBudgetAuth(r), beans.NewMonthDate(date))
		if err != nil {
			Error(w, err)
			return
		}

		responseCategories := make([]response.MonthCategory, len(month.Categories))
		for i, category := range month.Categories {
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
				Date:        month.Date,
				Budgetable:  month.Budgetable,
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
	return func(w http.ResponseWriter, r *http.Request) {
		var req request.UpdateMonth
		if err := decodeRequest(r, &req); err != nil {
			Error(w, err)
			return
		}

		monthID, err := beans.IDFromString(chi.URLParam(r, "monthID"))
		if err != nil {
			Error(w, err)
			return
		}

		if err := s.contracts.Month.Update(r.Context(), getBudgetAuth(r), monthID, req.Carryover); err != nil {
			Error(w, err)
			return
		}
	}
}

func (s *Server) handleMonthCategoryUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request.UpdateMonthCategory
		if err := decodeRequest(r, &req); err != nil {
			Error(w, err)
			return
		}

		monthID, err := beans.IDFromString(chi.URLParam(r, "monthID"))
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
