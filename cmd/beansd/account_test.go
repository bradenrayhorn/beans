package main_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanCreateAccount(t *testing.T) {
	ta := StartApplication(t)
	defer ta.Stop(t)

	user, session := ta.CreateUserAndSession(t)
	budget := ta.CreateBudget(t, "my budget", user)
	r := ta.PostRequest(t, fmt.Sprintf("api/v1/budgets/%s/accounts", budget.ID), &RequestOptions{SessionID: string(session.ID), Body: `{"name": "account1"}`})
	assert.Equal(t, http.StatusOK, r.StatusCode)
}

func TestCannotCreateAccountWithNoName(t *testing.T) {
	ta := StartApplication(t)
	defer ta.Stop(t)

	user, session := ta.CreateUserAndSession(t)
	budget := ta.CreateBudget(t, "my budget", user)
	r := ta.PostRequest(t, fmt.Sprintf("api/v1/budgets/%s/accounts", budget.ID), &RequestOptions{SessionID: string(session.ID), Body: `{"name": ""}`})
	assert.Equal(t, http.StatusUnprocessableEntity, r.StatusCode)
}

func TestCannotCreateAccountWithOthersBudget(t *testing.T) {
	ta := StartApplication(t)
	defer ta.Stop(t)

	_, session := ta.CreateUserAndSession(t)
	user2, _ := ta.CreateUserAndSession(t)
	budget := ta.CreateBudget(t, "my budget", user2)
	r := ta.PostRequest(t, fmt.Sprintf("api/v1/budgets/%s/accounts", budget.ID), &RequestOptions{SessionID: string(session.ID), Body: `{"name": "account1"}`})
	assert.Equal(t, http.StatusForbidden, r.StatusCode)
}

func TestCanGetAccounts(t *testing.T) {
	ta := StartApplication(t)
	defer ta.Stop(t)

	user, session := ta.CreateUserAndSession(t)
	budget := ta.CreateBudget(t, "my budget", user)
	r := ta.GetRequest(t, fmt.Sprintf("api/v1/budgets/%s/accounts", budget.ID), &RequestOptions{SessionID: string(session.ID)})
	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.JSONEq(t, `{"data":[]}`, r.Body)

	account := ta.CreateAccount(t, "account1", budget)
	r = ta.GetRequest(t, fmt.Sprintf("api/v1/budgets/%s/accounts", budget.ID), &RequestOptions{SessionID: string(session.ID)})
	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.JSONEq(t, fmt.Sprintf(`{"data":[{"name":"%s","id":"%s"}]}`, "account1", account.ID), r.Body)
}
