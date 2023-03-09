package http

import (
	"net/http"

	"github.com/bradenrayhorn/beans/beans"
)

type listCategoryResponse struct {
	ID   beans.ID   `json:"id"`
	Name beans.Name `json:"name"`
}

type categoryGroupResponse struct {
	ID         beans.ID               `json:"id"`
	Name       beans.Name             `json:"name"`
	IsIncome   bool                   `json:"is_income"`
	Categories []listCategoryResponse `json:"categories"`
}

func (s *Server) handleCategoryCreate() http.HandlerFunc {
	type request struct {
		GroupID beans.ID   `json:"group_id"`
		Name    beans.Name `json:"name"`
	}
	type response struct {
		ID beans.ID `json:"id"`
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

		jsonResponse(w, struct {
			Data response `json:"data"`
		}{
			Data: response{ID: category.ID},
		}, http.StatusOK)
	}
}

func (s *Server) handleCategoryGroupCreate() http.HandlerFunc {
	type request struct {
		Name beans.Name `json:"name"`
	}
	type response struct {
		ID beans.ID `json:"id"`
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

		jsonResponse(w, struct {
			Data response `json:"data"`
		}{
			Data: response{ID: group.ID},
		}, http.StatusOK)
	}
}

func (s *Server) handleCategoryGetAll() http.HandlerFunc {
	type response struct {
		Data []categoryGroupResponse `json:"data"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		groups, categories, err := s.categoryContract.GetAll(r.Context(), getBudgetAuth(r))
		if err != nil {
			Error(w, err)
			return
		}

		categoriesMap := make(map[string][]listCategoryResponse)
		for _, group := range groups {
			categoriesMap[group.ID.String()] = make([]listCategoryResponse, 0)
		}
		for _, category := range categories {
			groupID := category.GroupID.String()
			categoriesMap[groupID] = append(categoriesMap[groupID], listCategoryResponse{ID: category.ID, Name: category.Name})
		}

		res := response{Data: make([]categoryGroupResponse, len(groups))}
		for i, group := range groups {
			res.Data[i] = categoryGroupResponse{
				ID:         group.ID,
				Name:       group.Name,
				IsIncome:   group.IsIncome,
				Categories: categoriesMap[group.ID.String()],
			}
		}

		jsonResponse(w, res, http.StatusOK)
	}
}
