package httpadapter

import (
	"fmt"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/http/response"
	"github.com/bradenrayhorn/beans/server/specification"
)

func (a *httpAdapter) PayeeCreate(t *testing.T, ctx specification.Context, name beans.Name) (beans.ID, error) {
	r := a.Request(t, HTTPRequest{
		Method:  "POST",
		Path:    "/api/v1/payees",
		Body:    fmt.Sprintf(`{"name":"%s"}`, name),
		Context: ctx,
	})
	resp, err := MustParseResponse[response.CreatePayeeResponse](t, r.Response)
	if err != nil {
		return beans.ID{}, err
	}
	return resp.Data.ID, nil
}

func (a *httpAdapter) PayeeGetAll(t *testing.T, ctx specification.Context) ([]beans.Payee, error) {
	r := a.Request(t, HTTPRequest{
		Method:  "GET",
		Path:    "/api/v1/payees",
		Context: ctx,
	})
	resp, err := MustParseResponse[response.ListPayeesResponse](t, r.Response)
	if err != nil {
		return nil, err
	}

	return mapAll(resp.Data, mapPayee), nil
}

func (a *httpAdapter) PayeeGet(t *testing.T, ctx specification.Context, id beans.ID) (beans.Payee, error) {
	r := a.Request(t, HTTPRequest{
		Method:  "GET",
		Path:    fmt.Sprintf("/api/v1/payees/%s", id),
		Context: ctx,
	})
	resp, err := MustParseResponse[response.GetPayeeResponse](t, r.Response)
	if err != nil {
		return beans.Payee{}, err
	}

	return mapPayee(resp.Data), nil
}
