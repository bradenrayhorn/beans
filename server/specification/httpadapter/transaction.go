package httpadapter

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/http/request"
	"github.com/bradenrayhorn/beans/server/http/response"
	"github.com/bradenrayhorn/beans/server/specification"
	"github.com/stretchr/testify/require"
)

func mustEncode(t *testing.T, v any) string {
	bytes, err := json.Marshal(v)
	require.NoError(t, err, "could not encode: %v", v)
	return string(bytes)
}

func (a *HTTPAdapter) TransactionCreate(t *testing.T, ctx specification.Context, params beans.TransactionCreateParams) (beans.ID, error) {
	r := a.Request(t, HTTPRequest{
		Method: "POST",
		Path:   "/api/v1/transactions",
		Body: mustEncode(t, request.CreateTransactionRequest{
			AccountID:  params.AccountID,
			CategoryID: params.CategoryID,
			PayeeID:    params.PayeeID,
			Amount:     params.Amount,
			Date:       params.Date,
			Notes:      params.Notes,
		}),
		Context: ctx,
	})
	resp, err := MustParseResponse[response.CreateTransactionResponse](t, r.Response)
	if err != nil {
		return beans.ID{}, err
	}
	return resp.Data.ID, nil
}

func (a *HTTPAdapter) TransactionGet(t *testing.T, ctx specification.Context, id beans.ID) (beans.TransactionWithRelations, error) {
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
