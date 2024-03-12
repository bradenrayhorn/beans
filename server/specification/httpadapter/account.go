package httpadapter

import (
	"fmt"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/http/request"
	"github.com/bradenrayhorn/beans/server/http/response"
	"github.com/bradenrayhorn/beans/server/specification"
)

func (a *httpAdapter) AccountCreate(t *testing.T, ctx specification.Context, params beans.AccountCreate) (beans.ID, error) {
	r := a.Request(t, HTTPRequest{
		Method: "POST",
		Path:   "/api/v1/accounts",
		Body: mustEncode(t, request.CreateAccount{
			Name:      params.Name,
			OffBudget: params.OffBudget,
		}),
		Context: ctx,
	})
	resp, err := MustParseResponse[response.CreateAccountResponse](t, r.Response)
	if err != nil {
		return beans.ID{}, err
	}
	return resp.Data.ID, nil
}

func (a *httpAdapter) AccountList(t *testing.T, ctx specification.Context) ([]beans.AccountWithBalance, error) {
	r := a.Request(t, HTTPRequest{
		Method:  "GET",
		Path:    "/api/v1/accounts",
		Context: ctx,
	})
	resp, err := MustParseResponse[response.ListAccountResponse](t, r.Response)
	if err != nil {
		return nil, err
	}
	return mapAll(resp.Data, mapListAccount), nil
}

func (a *httpAdapter) AccountListTransactable(t *testing.T, ctx specification.Context) ([]beans.Account, error) {
	r := a.Request(t, HTTPRequest{
		Method:  "GET",
		Path:    "/api/v1/accounts/transactable",
		Context: ctx,
	})
	resp, err := MustParseResponse[response.GetTransactableAccounts](t, r.Response)
	if err != nil {
		return nil, err
	}
	return mapAll(resp.Data, mapAccount), nil
}

func (a *httpAdapter) AccountGet(t *testing.T, ctx specification.Context, id beans.ID) (beans.Account, error) {
	r := a.Request(t, HTTPRequest{
		Method:  "GET",
		Path:    fmt.Sprintf("/api/v1/accounts/%s", id),
		Context: ctx,
	})
	resp, err := MustParseResponse[response.GetAccountResponse](t, r.Response)
	if err != nil {
		return beans.Account{}, err
	}

	return mapAccount(resp.Data), nil
}
