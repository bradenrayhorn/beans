package httpadapter

import (
	"fmt"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/http/request"
	"github.com/bradenrayhorn/beans/server/http/response"
	"github.com/bradenrayhorn/beans/server/specification"
)

func (a *httpAdapter) MonthGetOrCreate(t *testing.T, ctx specification.Context, date beans.MonthDate) (beans.MonthWithDetails, error) {
	r := a.Request(t, HTTPRequest{
		Method:  "GET",
		Path:    fmt.Sprintf("/api/v1/months/%s", date),
		Context: ctx,
	})
	resp, err := MustParseResponse[response.GetMonthResponse](t, r.Response)
	if err != nil {
		return beans.MonthWithDetails{}, err
	}

	return mapMonthWithDetails(resp.Data), nil
}

func (a *httpAdapter) MonthUpdate(t *testing.T, ctx specification.Context, id beans.ID, carryover beans.Amount) error {
	r := a.Request(t, HTTPRequest{
		Method: "PUT",
		Path:   fmt.Sprintf("/api/v1/months/%s", id),
		Body: mustEncode(t, request.UpdateMonth{
			Carryover: carryover,
		}),
		Context: ctx,
	})
	if err := getErrorFromResponse(t, r.Response); err != nil {
		return err
	}

	return nil
}

func (a *httpAdapter) MonthSetCategoryAmount(t *testing.T, ctx specification.Context, id beans.ID, categoryID beans.ID, amount beans.Amount) error {
	r := a.Request(t, HTTPRequest{
		Method: "POST",
		Path:   fmt.Sprintf("/api/v1/months/%s/categories", id),
		Body: mustEncode(t, request.UpdateMonthCategory{
			CategoryID: categoryID,
			Amount:     amount,
		}),
		Context: ctx,
	})
	if err := getErrorFromResponse(t, r.Response); err != nil {
		t.Logf("got error: %v", err)
		return err
	}

	return nil
}
