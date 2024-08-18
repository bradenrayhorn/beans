package httpadapter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/http/response"
	"github.com/bradenrayhorn/beans/server/specification"
	"github.com/stretchr/testify/require"
)

type httpAdapter struct {
	// Base URL where the HTTP server is. Should not end in a slash.
	BaseURL string
}

var _ specification.Interactor = (*httpAdapter)(nil)

func New(baseURL string) specification.Interactor {
	return &httpAdapter{BaseURL: baseURL}
}

// Map HTTP status to bean code

var httpStatusToCode = map[int]string{
	http.StatusForbidden:           beans.EFORBIDDEN,
	http.StatusInternalServerError: beans.EINTERNAL,
	http.StatusUnprocessableEntity: beans.EINVALID,
	http.StatusNotFound:            beans.ENOTFOUND,
	http.StatusUnauthorized:        beans.EUNAUTHORIZED,
	http.StatusBadRequest:          beans.EUNPROCESSABLE,
}

// Request helpers

type HTTPRequest struct {
	Method  string
	Path    string
	Body    any
	Context specification.Context
}

type HTTPResponse struct {
	*http.Response
}

func (a *httpAdapter) Request(t *testing.T, req HTTPRequest) *HTTPResponse {
	// parse body
	var body io.Reader
	switch rawBody := req.Body.(type) {
	case string:
		body = bytes.NewReader([]byte(rawBody))
	case nil:
		body = nil
	default:
		body = nil
	}

	// create http request
	httpRequest, err := http.NewRequest(
		req.Method,
		a.BaseURL+req.Path,
		body,
	)
	require.Nil(t, err)

	// attach session id cookie
	if len(req.Context.SessionID) != 0 {
		httpRequest.Header.Add("Authorization", string(req.Context.SessionID))
	}

	// attach budget id header
	if !req.Context.BudgetID.Empty() {
		httpRequest.Header.Add("Budget-ID", req.Context.BudgetID.String())
	}

	client := &http.Client{}
	resp, err := client.Do(httpRequest)
	require.Nil(t, err)

	return &HTTPResponse{Response: resp}
}

func MustParseResponse[T any](t *testing.T, r *http.Response) (T, error) {
	var result T

	if err := getErrorFromResponse(t, r); err != nil {
		return result, err
	} else {
		defer func() { _ = r.Body.Close() }()
		require.Nil(t, json.NewDecoder(r.Body).Decode(&result))
		return result, nil
	}
}

func getErrorFromResponse(t *testing.T, r *http.Response) error {
	if r.StatusCode >= 200 && r.StatusCode <= 299 {
		return nil
	} else {
		// Else parse the error and transform into a beans error
		defer func() { _ = r.Body.Close() }()
		var resp response.ErrorResponse
		if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
			if errors.Is(err, io.EOF) {
				return fmt.Errorf("request error, no body. status code %d", r.StatusCode)
			}

			return fmt.Errorf("request error, unknown response. status code %d. error: %v", r.StatusCode, err)
		}
		code := beans.EINTERNAL
		if val, ok := httpStatusToCode[r.StatusCode]; ok {
			code = val
		}
		return beans.NewError(code, resp.Error)
	}
}

func mustEncode(t *testing.T, v any) string {
	bytes, err := json.Marshal(v)
	require.NoError(t, err, "could not encode: %v", v)
	return string(bytes)
}
