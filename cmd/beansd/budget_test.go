package main_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBudgets(t *testing.T) {
	ta := StartApplication(t)
	defer ta.Stop(t)

	t.Run("can create budget", func(t *testing.T) {
		_, session := ta.CreateUserAndSession(t)
		r := ta.PostRequest(t, "api/v1/budgets", &RequestOptions{SessionID: string(session.ID), Body: `{"name": "my budget"}`})
		assert.Equal(t, http.StatusOK, r.StatusCode)
	})

	t.Run("create requires name", func(t *testing.T) {
		_, session := ta.CreateUserAndSession(t)
		r := ta.PostRequest(t, "api/v1/budgets", &RequestOptions{SessionID: string(session.ID), Body: `{"name":""}`})
		assert.Equal(t, http.StatusUnprocessableEntity, r.StatusCode)
		assert.JSONEq(t, `{"error":"Budget name is required.","code":"invalid"}`, r.Body)
	})

	t.Run("get budgets", func(t *testing.T) {
		user, session := ta.CreateUserAndSession(t)
		r := ta.GetRequest(t, "api/v1/budgets", &RequestOptions{SessionID: string(session.ID)})
		assert.Equal(t, http.StatusOK, r.StatusCode)
		assert.JSONEq(t, `{"data":[]}`, r.Body)

		budget := ta.CreateBudget(t, "my budget", user)
		r = ta.GetRequest(t, "api/v1/budgets", &RequestOptions{SessionID: string(session.ID)})
		assert.Equal(t, http.StatusOK, r.StatusCode)
		assert.JSONEq(t, fmt.Sprintf(`{"data":[{"name":"%s","id":"%s"}]}`, "my budget", budget.ID), r.Body)
	})

	t.Run("can get budget", func(t *testing.T) {
		user, session := ta.CreateUserAndSession(t)
		budget := ta.CreateBudget(t, "my budget", user)
		month, err := ta.application.MonthRepository().GetLatest(context.Background(), budget.ID)
		require.Nil(t, err)
		r := ta.GetRequest(t, fmt.Sprintf("api/v1/budgets/%s", budget.ID), &RequestOptions{SessionID: string(session.ID)})

		assert.Equal(t, http.StatusOK, r.StatusCode)
		assert.JSONEq(t, fmt.Sprintf(`{"data":{"name":"%s","id":"%s","latest_month_id":"%s"}}`, "my budget", budget.ID, month.ID), r.Body)
	})
}
