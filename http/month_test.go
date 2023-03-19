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
		contract.SetCategoryAmountFunc.PushReturn(nil)

		categoryID := beans.NewBeansID()

		req := fmt.Sprintf(`{"category_id":"%s","amount":34}`, categoryID)
		options := &testutils.HTTPOptions{URLParams: map[string]string{"monthID": month.ID.String()}}
		res := testutils.HTTPWithOptions(t, sv.handleMonthCategoryUpdate(), options, user, budget, req, http.StatusOK)

		assert.Empty(t, res)

		params := contract.SetCategoryAmountFunc.History()[0]
		assert.Equal(t, budget.ID, params.Arg1.BudgetID())
		assert.Equal(t, month.ID, params.Arg2)
		assert.Equal(t, categoryID, params.Arg3)
		assert.Equal(t, beans.NewAmount(34, 0), params.Arg4)
	})

	t.Run("update month", func(t *testing.T) {
		contract.UpdateFunc.PushReturn(nil)

		req := `{"carryover":34}`
		options := &testutils.HTTPOptions{URLParams: map[string]string{"monthID": month.ID.String()}}
		res := testutils.HTTPWithOptions(t, sv.handleMonthUpdate(), options, user, budget, req, http.StatusOK)

		assert.Empty(t, res)

		params := contract.UpdateFunc.History()[0]
		assert.Equal(t, budget.ID, params.Arg1.BudgetID())
		assert.Equal(t, month.ID, params.Arg2)
		assert.Equal(t, beans.NewAmount(34, 0), params.Arg3)
	})

	t.Run("get", func(t *testing.T) {
		category := &beans.MonthCategory{ID: beans.NewBeansID(), CategoryID: beans.NewBeansID(), Amount: beans.NewAmount(5, 0), Activity: beans.NewAmount(4, 0), Available: beans.NewAmount(1, 0)}
		contract.GetOrCreateFunc.PushReturn(month, []*beans.MonthCategory{category}, beans.NewAmount(55, 0), nil)

		options := &testutils.HTTPOptions{URLParams: map[string]string{"date": "2022-05-01"}}
		res := testutils.HTTPWithOptions(t, sv.handleMonthGetOrCreate(), options, user, budget, nil, http.StatusOK)

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

		assert.JSONEq(t, expected, res)
	})

	t.Run("cannot get with invalid date", func(t *testing.T) {
		options := &testutils.HTTPOptions{URLParams: map[string]string{"date": "2022---33"}}
		res := testutils.HTTPWithOptions(t, sv.handleMonthGetOrCreate(), options, user, budget, nil, http.StatusUnprocessableEntity)

		expected := `{"error":"Invalid data provided","code":"invalid"}`
		assert.JSONEq(t, expected, res)
	})

	t.Run("cannot get no date", func(t *testing.T) {
		options := &testutils.HTTPOptions{URLParams: map[string]string{"date": ""}}
		res := testutils.HTTPWithOptions(t, sv.handleMonthGetOrCreate(), options, user, budget, nil, http.StatusUnprocessableEntity)

		expected := `{"error":"Invalid data provided","code":"invalid"}`
		assert.JSONEq(t, expected, res)
	})
}
