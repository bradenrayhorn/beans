package httpadapter

import (
	"fmt"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/http/response"
	"github.com/bradenrayhorn/beans/server/specification"
)

func (a *httpAdapter) CategoryGroupCreate(t *testing.T, ctx specification.Context, name beans.Name) (beans.ID, error) {
	r := a.Request(t, HTTPRequest{
		Method:  "POST",
		Path:    "/api/v1/categories/groups",
		Body:    fmt.Sprintf(`{"name":"%s"}`, name),
		Context: ctx,
	})
	resp, err := MustParseResponse[response.CreateCategoryGroupResponse](t, r.Response)
	if err != nil {
		return beans.ID{}, err
	}
	return resp.Data.ID, nil
}

func (a *httpAdapter) CategoryGroupGet(t *testing.T, ctx specification.Context, id beans.ID) (beans.CategoryGroup, error) {
	r := a.Request(t, HTTPRequest{
		Method:  "GET",
		Path:    fmt.Sprintf("/api/v1/categories/groups/%s", id),
		Context: ctx,
	})
	resp, err := MustParseResponse[response.GetCategoryGroupResponse](t, r.Response)
	if err != nil {
		return beans.CategoryGroup{}, err
	}

	return mapCategoryGroup(resp.Data), nil
}

func (a *httpAdapter) CategoryCreate(t *testing.T, ctx specification.Context, groupID beans.ID, name beans.Name) (beans.ID, error) {
	r := a.Request(t, HTTPRequest{
		Method:  "POST",
		Path:    "/api/v1/categories",
		Body:    fmt.Sprintf(`{"name":"%s","group_id":"%s"}`, name, groupID),
		Context: ctx,
	})
	resp, err := MustParseResponse[response.CreateCategoryResponse](t, r.Response)
	if err != nil {
		return beans.ID{}, err
	}
	return resp.Data.ID, nil
}

func (a *httpAdapter) CategoryGet(t *testing.T, ctx specification.Context, id beans.ID) (beans.Category, error) {
	r := a.Request(t, HTTPRequest{
		Method:  "GET",
		Path:    fmt.Sprintf("/api/v1/categories/%s", id),
		Context: ctx,
	})
	resp, err := MustParseResponse[response.GetCategoryResponse](t, r.Response)
	if err != nil {
		return beans.Category{}, err
	}

	return mapCategory(resp.Data), nil
}
