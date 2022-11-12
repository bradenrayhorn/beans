package main_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestMonths(t *testing.T) {
	ta := StartApplication(t)
	defer ta.Stop(t)

	t.Run("can get month", func(t *testing.T) {
		user, session := ta.CreateUserAndSession(t)
		budget := ta.CreateBudget(t, "my budget", user)
		month := ta.CreateMonth(t, budget, testutils.NewDate(t, "2022-05-01"))
		categoryGroup := ta.CreateCategoryGroup(t, budget, "Bills")
		category := ta.CreateCategory(t, budget, categoryGroup, "Electric")
		monthCategory := ta.CreateMonthCategory(t, month, category, beans.NewAmount(5, -1))

		r := ta.GetRequest(t, fmt.Sprintf("api/v1/months/%s", month.ID), newOptions(session, budget))
		assert.Equal(t, http.StatusOK, r.StatusCode)
		assert.JSONEq(t, fmt.Sprintf(`{"data":{
			"id": "%s",
			"date": "2022-05-01",
			"categories": [
				{
					"id": "%s",
					"category_id": "%s",
					"assigned": {
						"coefficient": 5,
						"exponent": -1
					}
				}
			]
		}}`, month.ID, monthCategory.ID, category.ID), r.Body)
	})
}
