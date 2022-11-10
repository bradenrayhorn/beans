package http

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/mocks"
	"github.com/bradenrayhorn/beans/internal/testutils"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateBudget(t *testing.T) {
	budgetService := new(mocks.BudgetService)
	sv := &Server{budgetService: budgetService}
	user := &beans.User{ID: beans.UserID(beans.NewBeansID())}
	budget := &beans.Budget{ID: beans.NewBeansID(), Name: "Budget1"}

	t.Run("create calls service and returns response", func(t *testing.T) {
		call := budgetService.On("CreateBudget", mock.Anything, budget.Name, user.ID).Return(budget, nil)
		defer call.Unset()

		req := fmt.Sprintf(`{
      "name": "%s"
    }`, "Budget1")

		resp := testutils.HTTP(t, sv.handleBudgetCreate(), user, nil, req, http.StatusOK)
		assert.JSONEq(t, resp, fmt.Sprintf(`{"data":{
      "id": "%s",
      "name": "Budget1"
    }}`, budget.ID))
	})
}

func TestGetSingleBudget(t *testing.T) {
	monthService := new(mocks.MonthService)
	budgetRepository := new(mocks.BudgetRepository)
	sv := &Server{monthService: monthService, budgetRepository: budgetRepository}
	user := &beans.User{ID: beans.UserID(beans.NewBeansID())}
	budget := &beans.Budget{ID: beans.NewBeansID(), Name: "Budget1", UserIDs: []beans.UserID{user.ID}}
	month := &beans.Month{ID: beans.NewBeansID(), BudgetID: budget.ID, Date: beans.NewDate(time.Now())}

	t.Run("can get budget", func(t *testing.T) {
		call := budgetRepository.On("Get", mock.Anything, budget.ID).Return(budget, nil)
		defer call.Unset()
		call = monthService.On("GetOrCreate", mock.Anything, budget.ID, mock.Anything).Return(month, nil)
		defer call.Unset()

		withContext := func(ctx context.Context) context.Context {
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("budgetID", budget.ID.String())
			return context.WithValue(ctx, chi.RouteCtxKey, rctx)
		}
		resp := testutils.HTTPWithContext(t, sv.handleBudgetGet(), withContext, user, budget, nil, http.StatusOK)
		assert.JSONEq(t, resp, fmt.Sprintf(`{"data":{
      "id": "%s",
      "name": "%s",
      "default_month_id": "%s"
    }}`, budget.ID, budget.Name, month.ID))
	})
}

func TestGetAllBudgets(t *testing.T) {
	budgetRepository := new(mocks.BudgetRepository)
	sv := &Server{budgetRepository: budgetRepository}
	user := &beans.User{ID: beans.UserID(beans.NewBeansID())}
	budget1 := &beans.Budget{ID: beans.NewBeansID(), Name: "Budget1", UserIDs: []beans.UserID{user.ID}}
	budget2 := &beans.Budget{ID: beans.NewBeansID(), Name: "Budget2", UserIDs: []beans.UserID{user.ID}}

	t.Run("can get all budgets", func(t *testing.T) {
		call := budgetRepository.On("GetBudgetsForUser", mock.Anything, user.ID).Return([]*beans.Budget{budget1, budget2}, nil)
		defer call.Unset()

		resp := testutils.HTTP(t, sv.handleBudgetGetAll(), user, nil, nil, http.StatusOK)
		assert.JSONEq(t, resp, fmt.Sprintf(`{"data":[
      {
        "id": "%s",
        "name": "%s"
      },
      {
        "id": "%s",
        "name": "%s"
      }
    ]}`, budget1.ID, budget1.Name, budget2.ID, budget2.Name))
	})
}
