package main_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanCreateBudget(t *testing.T) {
	ta := StartApplication(t)
	defer ta.Stop(t)

	_, session := ta.CreateUserAndSession(t)
	r := ta.PostRequest(t, "api/v1/budgets", &RequestOptions{SessionID: string(session.ID), Body: `{"name": "my budget"}`})
	assert.Equal(t, http.StatusOK, r.StatusCode)
}

func TestBudgetCreationRequiresName(t *testing.T) {
	ta := StartApplication(t)
	defer ta.Stop(t)

	_, session := ta.CreateUserAndSession(t)
	r := ta.PostRequest(t, "api/v1/budgets", &RequestOptions{SessionID: string(session.ID), Body: `{"name":""}`})
	assert.Equal(t, http.StatusUnprocessableEntity, r.StatusCode)
	assert.JSONEq(t, `{"error":"Budget name is required.","code":"invalid"}`, r.Body)
}

func TestGetBudgets(t *testing.T) {
	ta := StartApplication(t)
	defer ta.Stop(t)

	user, session := ta.CreateUserAndSession(t)
	r := ta.GetRequest(t, "api/v1/budgets", &RequestOptions{SessionID: string(session.ID)})
	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.JSONEq(t, `{"data":[]}`, r.Body)

	budget := ta.CreateBudget(t, "my budget", user)
	r = ta.GetRequest(t, "api/v1/budgets", &RequestOptions{SessionID: string(session.ID)})
	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.JSONEq(t, fmt.Sprintf(`{"data":[{"name":"%s","id":"%s"}]}`, "my budget", budget.ID), r.Body)

}
