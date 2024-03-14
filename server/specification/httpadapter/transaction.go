package httpadapter

import (
	"fmt"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/http/request"
	"github.com/bradenrayhorn/beans/server/http/response"
	"github.com/bradenrayhorn/beans/server/specification"
)

func (a *httpAdapter) TransactionCreate(t *testing.T, ctx specification.Context, params beans.TransactionCreateParams) (beans.ID, error) {
	r := a.Request(t, HTTPRequest{
		Method: "POST",
		Path:   "/api/v1/transactions",
		Body: mustEncode(t, request.CreateTransaction{
			AccountID:         params.AccountID,
			CategoryID:        params.CategoryID,
			PayeeID:           params.PayeeID,
			Amount:            params.Amount,
			Date:              params.Date,
			Notes:             params.Notes,
			TransferAccountID: params.TransferAccountID,
		}),
		Context: ctx,
	})
	resp, err := MustParseResponse[response.CreateTransactionResponse](t, r.Response)
	if err != nil {
		return beans.ID{}, err
	}
	return resp.Data.ID, nil
}

func (a *httpAdapter) TransactionGet(t *testing.T, ctx specification.Context, id beans.ID) (beans.TransactionWithRelations, error) {
	r := a.Request(t, HTTPRequest{
		Method:  "GET",
		Path:    fmt.Sprintf("/api/v1/transactions/%s", id),
		Context: ctx,
	})
	resp, err := MustParseResponse[response.GetTransactionResponse](t, r.Response)
	if err != nil {
		return beans.TransactionWithRelations{}, err
	}

	return mapTransactionWithRelations(resp.Data), nil
}

func (a *httpAdapter) TransactionUpdate(t *testing.T, ctx specification.Context, params beans.TransactionUpdateParams) error {
	r := a.Request(t, HTTPRequest{
		Method: "PUT",
		Path:   fmt.Sprintf("/api/v1/transactions/%s", params.ID),
		Body: mustEncode(t, request.UpdateTransaction{
			AccountID:  params.AccountID,
			CategoryID: params.CategoryID,
			PayeeID:    params.PayeeID,
			Amount:     params.Amount,
			Date:       params.Date,
			Notes:      params.Notes,
		}),
		Context: ctx,
	})
	return getErrorFromResponse(t, r.Response)
}

func (a *httpAdapter) TransactionDelete(t *testing.T, ctx specification.Context, ids []beans.ID) error {
	r := a.Request(t, HTTPRequest{
		Method:  "POST",
		Path:    "/api/v1/transactions/delete",
		Body:    mustEncode(t, request.DeleteTransaction{IDs: ids}),
		Context: ctx,
	})
	return getErrorFromResponse(t, r.Response)
}

func (a *httpAdapter) TransactionGetAll(t *testing.T, ctx specification.Context) ([]beans.TransactionWithRelations, error) {
	r := a.Request(t, HTTPRequest{
		Method:  "GET",
		Path:    "/api/v1/transactions",
		Context: ctx,
	})
	resp, err := MustParseResponse[response.ListTransactionsResponse](t, r.Response)
	if err != nil {
		return nil, err
	}

	return mapAll(resp.Data, mapTransactionWithRelations), nil
}
