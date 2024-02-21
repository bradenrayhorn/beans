package http

import (
	"net/http"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/http/response"
	"github.com/go-chi/chi/v5"
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

		category, err := s.contracts.Category.CreateCategory(r.Context(), getBudgetAuth(r), req.GroupID, req.Name)
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

		group, err := s.contracts.Category.CreateGroup(r.Context(), getBudgetAuth(r), req.Name)
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
		groups, err := s.contracts.Category.GetAll(r.Context(), getBudgetAuth(r))
		if err != nil {
			Error(w, err)
			return
		}

		res := response.GetCategoriesResponse{Data: make([]response.CategoryGroup, len(groups))}
		for i, group := range groups {
			categories := []response.AssociatedCategory{}

			for _, category := range group.Categories {
				categories = append(categories, response.AssociatedCategory{ID: category.ID, Name: category.Name})
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

func (s *Server) handleCategoryGetCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := beans.IDFromString(chi.URLParam(r, "categoryID"))
		if err != nil {
			Error(w, beans.WrapError(err, beans.ErrorNotFound))
			return
		}

		category, err := s.contracts.Category.GetCategory(r.Context(), getBudgetAuth(r), id)
		if err != nil {
			Error(w, err)
			return
		}

		jsonResponse(w, response.GetCategoryResponse{
			Data: response.Category{
				ID:      category.ID,
				Name:    category.Name,
				GroupID: category.GroupID,
			},
		}, http.StatusOK)
	}
}

func (s *Server) handleCategoryGetCategoryGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := beans.IDFromString(chi.URLParam(r, "categoryGroupID"))
		if err != nil {
			Error(w, beans.WrapError(err, beans.ErrorNotFound))
			return
		}

		group, err := s.contracts.Category.GetGroup(r.Context(), getBudgetAuth(r), id)
		if err != nil {
			Error(w, err)
			return
		}

		categories := []response.AssociatedCategory{}
		for _, category := range group.Categories {
			categories = append(categories, response.AssociatedCategory{ID: category.ID, Name: category.Name})
		}

		jsonResponse(w, response.GetCategoryGroupResponse{
			Data: response.CategoryGroup{
				ID:         group.ID,
				Name:       group.Name,
				IsIncome:   group.IsIncome,
				Categories: categories,
			},
		}, http.StatusOK)
	}
}
