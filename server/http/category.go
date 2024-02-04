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
		groups, categories, err := s.categoryContract.GetAll(r.Context(), getBudgetAuth(r))
		if err != nil {
			Error(w, err)
			return
		}

		categoriesMap := make(map[string][]response.Category)
		for _, group := range groups {
			categoriesMap[group.ID.String()] = make([]response.Category, 0)
		}
		for _, category := range categories {
			groupID := category.GroupID.String()
			categoriesMap[groupID] = append(categoriesMap[groupID], response.Category{ID: category.ID, Name: category.Name})
		}

		res := response.GetCategoriesResponse{Data: make([]response.CategoryGroup, len(groups))}
		for i, group := range groups {
			res.Data[i] = response.CategoryGroup{
				ID:         group.ID,
				Name:       group.Name,
				IsIncome:   group.IsIncome,
				Categories: categoriesMap[group.ID.String()],
			}
		}

		jsonResponse(w, res, http.StatusOK)
	}
}
