package http

import (
	"net/http"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/http/response"
)

func (s *Server) handleCategoryCreate() http.HandlerFunc {
	type request struct {
		GroupID beans.ID   `json:"group_id"`
		Name    beans.Name `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := decodeRequest(r, &req); err != nil {
			Error(w, err)
			return
		}

		category, err := s.categoryContract.CreateCategory(r.Context(), getBudgetAuth(r), req.GroupID, req.Name)
		if err != nil {
			Error(w, err)
			return
		}

		jsonResponse(w, response.CreateCategoryResponse{
			Data: response.ID{ID: category.ID},
		}, http.StatusOK)
	}
}

func (s *Server) handleCategoryGroupCreate() http.HandlerFunc {
	type request struct {
		Name beans.Name `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := decodeRequest(r, &req); err != nil {
			Error(w, err)
			return
		}

		group, err := s.categoryContract.CreateGroup(r.Context(), getBudgetAuth(r), req.Name)
		if err != nil {
			Error(w, err)
			return
		}

		jsonResponse(w, response.CreateCategoryGroupResponse{
			Data: response.ID{ID: group.ID},
		}, http.StatusOK)
	}
}

func (s *Server) handleCategoryGetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groups, err := s.categoryContract.GetAll(r.Context(), getBudgetAuth(r))
		if err != nil {
			Error(w, err)
			return
		}

		res := response.GetCategoriesResponse{Data: make([]response.CategoryGroup, len(groups))}
		for i, group := range groups {
			categories := []response.Category{}

			for _, category := range group.Categories {
				categories = append(categories, response.Category{ID: category.ID, Name: category.Name})
			}

			res.Data[i] = response.CategoryGroup{
				ID:         group.ID,
				Name:       group.Name,
				IsIncome:   group.IsIncome,
				Categories: categories,
			}
		}

		jsonResponse(w, res, http.StatusOK)
	}
}
