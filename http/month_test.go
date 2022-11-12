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

func TestGetMonth(t *testing.T) {
	repository := new(mocks.MonthRepository)
	sv := &Server{monthRepository: repository}

	user := &beans.User{ID: beans.UserID(beans.NewBeansID())}
	budget := &beans.Budget{ID: beans.NewBeansID(), Name: "Budget"}
	budget2 := &beans.Budget{ID: beans.NewBeansID(), Name: "Budget2"}

	date := beans.NewDate(time.Date(2022, 05, 26, 0, 0, 0, 0, time.UTC))
	month := &beans.Month{ID: beans.NewBeansID(), BudgetID: budget.ID, Date: date}

	repository.On("Get", mock.Anything, month.ID).Return(month, nil)

	t.Run("can get month", func(t *testing.T) {
		resp := testutils.HTTPWithContext(t, sv.handleMonthGet(), withID(month.ID), user, budget, nil, http.StatusOK)
		assert.JSONEq(t, resp, fmt.Sprintf(`{"data":{
			"id": "%s",
			"date": "2022-05-26"
		}}`, month.ID))
	})

	t.Run("cannot get month for wrong budget", func(t *testing.T) {
		testutils.HTTPWithContext(t, sv.handleMonthGet(), withID(month.ID), user, budget2, nil, http.StatusNotFound)
	})
}

func withID(id beans.ID) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", id.String())
		return context.WithValue(ctx, chi.RouteCtxKey, rctx)
	}
}
