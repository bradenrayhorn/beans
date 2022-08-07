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

func TestCanGetBudget(t *testing.T) {
	ta := StartApplication(t)
	defer ta.Stop(t)

	user, session := ta.CreateUserAndSession(t)
	budget := ta.CreateBudget(t, "my budget", user)
	r := ta.GetRequest(t, fmt.Sprintf("api/v1/budgets/%s", budget.ID), &RequestOptions{SessionID: string(session.ID)})
	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.JSONEq(t, fmt.Sprintf(`{"data":{"name":"%s","id":"%s"}}`, "my budget", budget.ID), r.Body)
}

func TestCannotGetOtherUsersBudget(t *testing.T) {
	ta := StartApplication(t)
	defer ta.Stop(t)

	_, session := ta.CreateUserAndSession(t)
	user2, _ := ta.CreateUserAndSession(t)
	budget := ta.CreateBudget(t, "my budget", user2)
	r := ta.GetRequest(t, fmt.Sprintf("api/v1/budgets/%s", budget.ID), &RequestOptions{SessionID: string(session.ID)})
	assert.Equal(t, http.StatusForbidden, r.StatusCode)
}

func TestCannotGetNonExistantBudget(t *testing.T) {
	ta := StartApplication(t)
	defer ta.Stop(t)

	_, session := ta.CreateUserAndSession(t)
	r := ta.GetRequest(t, fmt.Sprintf("api/v1/budgets/%s", "bad-id"), &RequestOptions{SessionID: string(session.ID)})
	assert.Equal(t, http.StatusNotFound, r.StatusCode)
}
