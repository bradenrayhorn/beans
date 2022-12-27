package main_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccounts(t *testing.T) {
	ta := StartApplication(t)
	defer ta.Stop(t)

	t.Run("can create account", func(t *testing.T) {
		user, session := ta.CreateUserAndSession(t)
		budget := ta.CreateBudget(t, "my budget", user)

		r := ta.PostRequest(t, "api/v1/accounts", newOptionsWithBody(session, budget, `{"name":"account1"}`))
		assert.Equal(t, http.StatusOK, r.StatusCode)

		accounts, err := ta.application.AccountRepository().GetForBudget(context.Background(), budget.ID)
		require.Nil(t, err)
		assert.Len(t, accounts, 1)
		assert.Equal(t, accounts[0].Name, beans.Name("account1"))

		assert.JSONEq(t, fmt.Sprintf(`{"data":{"id":"%s"}}`, accounts[0].ID), r.Body)
	})

	t.Run("can get accounts", func(t *testing.T) {
		user, session := ta.CreateUserAndSession(t)
		budget := ta.CreateBudget(t, "my budget", user)
		r := ta.GetRequest(t, "api/v1/accounts", newOptions(session, budget))
		assert.Equal(t, http.StatusOK, r.StatusCode)
		assert.JSONEq(t, `{"data":[]}`, r.Body)

		account := ta.CreateAccount(t, "account1", budget)
		r = ta.GetRequest(t, "api/v1/accounts", newOptions(session, budget))
		assert.Equal(t, http.StatusOK, r.StatusCode)
		assert.JSONEq(t, fmt.Sprintf(`{"data":[{"name":"%s","id":"%s"}]}`, "account1", account.ID), r.Body)
	})
}
