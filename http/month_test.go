package http

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/mocks"
	"github.com/bradenrayhorn/beans/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestMonth(t *testing.T) {
	contract := mocks.NewMockMonthContract()
	sv := Server{monthContract: contract}

	user := &beans.User{ID: beans.NewBeansID()}
	budget := &beans.Budget{ID: beans.NewBeansID(), Name: "Budget1"}
	month := &beans.Month{ID: beans.NewBeansID(), BudgetID: budget.ID, Date: testutils.NewMonthDate(t, "2022-05-01")}

	t.Run("create month", func(t *testing.T) {
		contract.CreateMonthFunc.PushReturn(month, nil)

		req := `{"date":"2022-05-01"}`
		res := testutils.HTTP(t, sv.handleMonthCreate(), user, budget, req, http.StatusOK)

		expected := fmt.Sprintf(`{"data":{"month_id":"%s"}}`, month.ID)
		assert.JSONEq(t, expected, res)

		params := contract.CreateMonthFunc.History()[0]
		assert.Equal(t, budget.ID, params.Arg1)
		assert.Equal(t, month.Date, params.Arg2)
	})

	t.Run("update month category", func(t *testing.T) {
		contract.SetCategoryAmountFunc.PushReturn(nil)

		categoryID := beans.NewBeansID()

		req := fmt.Sprintf(`{"category_id":"%s","amount":34}`, categoryID)
		options := &testutils.HTTPOptions{ContextValues: map[string]any{"month": month}}
		res := testutils.HTTPWithOptions(t, sv.handleMonthCategoryUpdate(), options, user, budget, req, http.StatusOK)

		assert.Empty(t, res)

		params := contract.SetCategoryAmountFunc.History()[0]
		assert.Equal(t, month.ID, params.Arg1)
		assert.Equal(t, categoryID, params.Arg2)
		assert.Equal(t, beans.NewAmount(34, 0), params.Arg3)
	})

	t.Run("get", func(t *testing.T) {
		category := &beans.MonthCategory{ID: beans.NewBeansID(), CategoryID: beans.NewBeansID(), Amount: beans.NewAmount(5, 0), Spent: beans.NewAmount(4, 0)}
		contract.GetFunc.PushReturn(month, []*beans.MonthCategory{category}, nil)

		options := &testutils.HTTPOptions{ContextValues: map[string]any{"month": month}}
		res := testutils.HTTPWithOptions(t, sv.handleMonthGet(), options, user, budget, nil, http.StatusOK)

		expected := fmt.Sprintf(`{"data": {
			"id": "%s",
			"date": "2022-05-01",
			"categories": [
				{
					"id": "%s",
					"category_id": "%s",
					"assigned": {
						"coefficient": 5,
						"exponent": 0
					},
					"spent": {
						"coefficient": 4,
						"exponent": 0
					}
				}
			]
		}}`, month.ID, category.ID, category.CategoryID)

		assert.JSONEq(t, expected, res)
	})
}
