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
		CategoryID beans.ID     `json:"category_id"`
	}
	type responseMonth struct {
		ID         beans.ID           `json:"id"`
		Date       string             `json:"date"`
		Categories []responseCategory `json:"categories"`
		// todo - return available $ to assign
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

		month, err := s.monthService.Get(r.Context(), monthID, getBudget(r).ID)
		if err != nil {
			Error(w, err)
			return
		}

		categories, err := s.monthCategoryRepository.GetForMonth(r.Context(), monthID)
		if err != nil {
			Error(w, err)
			return
		}

		responseCategories := make([]responseCategory, len(categories))
		for i, category := range categories {
			responseCategories[i] = responseCategory{
				ID:         category.ID,
				Assigned:   category.Amount,
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
