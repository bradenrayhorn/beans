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
	r := ta.PostRequest(t, "api/v1/accounts", newOptionsWithBody(session, budget, `{"name":"account1"}`))
	assert.Equal(t, http.StatusOK, r.StatusCode)
}

func TestCannotCreateAccountWithNoName(t *testing.T) {
	ta := StartApplication(t)
	defer ta.Stop(t)

	user, session := ta.CreateUserAndSession(t)
	budget := ta.CreateBudget(t, "my budget", user)
	r := ta.PostRequest(t, "api/v1/accounts", newOptionsWithBody(session, budget, `{"name":""}`))
	assert.Equal(t, http.StatusUnprocessableEntity, r.StatusCode)
}

func TestCannotCreateAccountWithOthersBudget(t *testing.T) {
	ta := StartApplication(t)
	defer ta.Stop(t)

	_, session := ta.CreateUserAndSession(t)
	user2, _ := ta.CreateUserAndSession(t)
	budget := ta.CreateBudget(t, "my budget", user2)
	r := ta.PostRequest(t, "api/v1/accounts", newOptionsWithBody(session, budget, `{"name":"account1"}`))
	assert.Equal(t, http.StatusForbidden, r.StatusCode)
}

func TestCanGetAccounts(t *testing.T) {
	ta := StartApplication(t)
	defer ta.Stop(t)

	user, session := ta.CreateUserAndSession(t)
	budget := ta.CreateBudget(t, "my budget", user)
	r := ta.GetRequest(t, "api/v1/accounts", newOptions(session, budget))
	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.JSONEq(t, `{"data":[]}`, r.Body)

	account := ta.CreateAccount(t, "account1", budget)
	r = ta.GetRequest(t, "api/v1/accounts", newOptions(session, budget))
	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.JSONEq(t, fmt.Sprintf(`{"data":[{"name":"%s","id":"%s"}]}`, "account1", account.ID), r.Body)
}
