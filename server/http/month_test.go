package http_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestMonth(t *testing.T) {
	user := &beans.User{ID: beans.NewBeansID()}
	budget := &beans.Budget{ID: beans.NewBeansID(), Name: "Budget1", UserIDs: []beans.ID{user.ID}}
	month := &beans.Month{
		ID:          beans.NewBeansID(),
		BudgetID:    budget.ID,
		Date:        testutils.NewMonthDate(t, "2022-05-01"),
		Carryover:   beans.NewAmount(5, 0),
		Income:      beans.NewAmount(6, 0),
		Assigned:    beans.NewAmount(7, 0),
		CarriedOver: beans.NewAmount(8, 0),
	}

	t.Run("update month category", func(t *testing.T) {
		test := newHttpTest(t)
		defer test.Stop(t)

		test.monthContract.SetCategoryAmountFunc.PushReturn(nil)

		categoryID := beans.NewBeansID()

		res := test.DoRequest(t, HTTPRequest{
			method: "POST",
			path:   fmt.Sprintf("/api/v1/months/%s/categories", month.ID),
			body:   fmt.Sprintf(`{"category_id":"%s","amount":34}`, categoryID),
			user:   user,
			budget: budget,
		})

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Empty(t, res.body)

		params := test.monthContract.SetCategoryAmountFunc.History()[0]
		assert.Equal(t, budget.ID, params.Arg1.BudgetID())
		assert.Equal(t, month.ID, params.Arg2)
		assert.Equal(t, categoryID, params.Arg3)
		assert.Equal(t, beans.NewAmount(34, 0), params.Arg4)
	})

	t.Run("update month", func(t *testing.T) {
		test := newHttpTest(t)
		defer test.Stop(t)

		test.monthContract.UpdateFunc.PushReturn(nil)

		res := test.DoRequest(t, HTTPRequest{
			method: "PUT",
			path:   fmt.Sprintf("/api/v1/months/%s", month.ID),
			body:   `{"carryover":34}`,
			user:   user,
			budget: budget,
		})

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Empty(t, res.body)

		params := test.monthContract.UpdateFunc.History()[0]
		assert.Equal(t, budget.ID, params.Arg1.BudgetID())
		assert.Equal(t, month.ID, params.Arg2)
		assert.Equal(t, beans.NewAmount(34, 0), params.Arg3)
	})

	t.Run("get", func(t *testing.T) {
		test := newHttpTest(t)
		defer test.Stop(t)

		category := &beans.MonthCategory{ID: beans.NewBeansID(), CategoryID: beans.NewBeansID(), Amount: beans.NewAmount(5, 0), Activity: beans.NewAmount(4, 0), Available: beans.NewAmount(1, 0)}

		test.monthContract.GetOrCreateFunc.PushReturn(month, []*beans.MonthCategory{category}, beans.NewAmount(55, 0), nil)

		res := test.DoRequest(t, HTTPRequest{
			method: "GET",
			path:   "/api/v1/months/2022-05-01",
			user:   user,
			budget: budget,
		})

		expected := fmt.Sprintf(`{"data": {
			"id": "%s",
			"date": "2022-05-01",
			"budgetable": {
				"coefficient": 55,
				"exponent": 0
			},
			"carryover": {
				"coefficient": 5,
				"exponent": 0
			},
			"income": {
				"coefficient": 6,
				"exponent": 0
			},
			"assigned": {
				"coefficient": 7,
				"exponent": 0
			},
			"carried_over": {
				"coefficient": 8,
				"exponent": 0
			},
			"categories": [
				{
					"id": "%s",
					"category_id": "%s",
					"assigned": {
						"coefficient": 5,
						"exponent": 0
					},
					"activity": {
						"coefficient": 4,
						"exponent": 0
					},
					"available": {
						"coefficient": 1,
						"exponent": 0
					}
				}
			]
		}}`, month.ID, category.ID, category.CategoryID)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.JSONEq(t, expected, res.body)
	})

	t.Run("cannot get with invalid date", func(t *testing.T) {
		test := newHttpTest(t)
		defer test.Stop(t)

		res := test.DoRequest(t, HTTPRequest{
			method: "GET",
			path:   "/api/v1/months/2022---33",
			user:   user,
			budget: budget,
		})

		expected := `{"error":"Invalid data provided","code":"invalid"}`
		assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)
		assert.JSONEq(t, expected, res.body)
	})
}
