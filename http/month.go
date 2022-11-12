package http

import (
	"net/http"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/go-chi/chi/v5"
)

func (s *Server) handleMonthGet() http.HandlerFunc {
	type responseMonth struct {
		ID   beans.ID `json:"id"`
		Date string   `json:"date"`
		// todo - return month_categories and available $ to assign
	}
	type response struct {
		Data responseMonth `json:"data"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		monthID, err := beans.BeansIDFromString(chi.URLParam(r, "id"))
		if err != nil {
			Error(w, err)
			return
		}

		month, err := s.monthRepository.Get(r.Context(), monthID)
		if err != nil {
			Error(w, err)
			return
		}

		if month.BudgetID != getBudget(r).ID {
			Error(w, beans.ErrorNotFound)
			return
		}

		jsonResponse(w, response{
			Data: responseMonth{
				ID:   month.ID,
				Date: month.Date.String(),
			},
		}, http.StatusOK)
	}
}
