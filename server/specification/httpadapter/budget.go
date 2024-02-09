package httpadapter

import (
	"fmt"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/http/response"
	"github.com/bradenrayhorn/beans/server/specification"
)

func (a *httpAdapter) BudgetCreate(t *testing.T, ctx specification.Context, name beans.Name) (beans.ID, error) {
	r := a.Request(t, HTTPRequest{
		Method:  "POST",
		Path:    "/api/v1/budgets",
		Body:    fmt.Sprintf(`{"name":"%s"}`, name),
		Context: ctx,
	})
	resp, err := MustParseResponse[response.CreateBudgetResponse](t, r.Response)
	if err != nil {
		return beans.ID{}, err
	}
	return resp.Data.ID, nil
}

func (a *httpAdapter) BudgetGet(t *testing.T, ctx specification.Context, id beans.ID) (beans.Budget, error) {
	r := a.Request(t, HTTPRequest{
		Method:  "GET",
		Path:    fmt.Sprintf("/api/v1/budgets/%s", id),
		Context: ctx,
	})
	resp, err := MustParseResponse[response.GetBudgetResponse](t, r.Response)
	if err != nil {
		return beans.Budget{}, err
	}

	return mapBudget(resp.Data), nil
}

func (a *httpAdapter) BudgetGetAll(t *testing.T, ctx specification.Context) ([]beans.Budget, error) {
	r := a.Request(t, HTTPRequest{
		Method:  "GET",
		Path:    "/api/v1/budgets",
		Context: ctx,
	})
	resp, err := MustParseResponse[response.ListBudgetsResponse](t, r.Response)
	if err != nil {
		return nil, err
	}

	return mapAll(resp.Data, mapBudget), nil
}
