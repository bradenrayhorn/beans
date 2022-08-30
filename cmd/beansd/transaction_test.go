package main_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanCreateTransaction(t *testing.T) {
	ta := StartApplication(t)
	defer ta.Stop(t)

	user, session := ta.CreateUserAndSession(t)
	budget := ta.CreateBudget(t, "my budget", user)
	account := ta.CreateAccount(t, "my account", budget)
	r := ta.PostRequest(t, "api/v1/transactions", newOptionsWithBody(session, budget, fmt.Sprintf(`{"account_id":"%s","amount":55,"date":"2022-08-29","notes":"good"}`, account.ID)))
	assert.Equal(t, http.StatusOK, r.StatusCode)
}
