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
	"github.com/bradenrayhorn/beans/server/internal/testutils"
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
		httpRequest.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: string(req.Context.SessionID),
		})
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
			if errors.Is(io.EOF, err) {
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

// userAndBudget

type userAndBudget struct {
	t         *testing.T
	sessionID beans.SessionID
	budget    beans.Budget
	context   specification.Context

	adapter *httpAdapter
}

var _ specification.TestUserAndBudget = (*userAndBudget)(nil)

func (u *userAndBudget) Ctx() specification.Context {
	return u.context
}

func (u *userAndBudget) Budget() beans.Budget {
	return u.budget
}

func (u *userAndBudget) Account(opt specification.AccountOpts) beans.Account {
	name := beans.Name(beans.NewBeansID().String())

	id, err := u.adapter.AccountCreate(u.t, u.context, name)
	require.NoError(u.t, err)
	account, err := u.adapter.AccountGet(u.t, u.context, id)
	require.NoError(u.t, err)

	return account
}

func (u *userAndBudget) CategoryGroup(opt specification.CategoryGroupOpts) beans.CategoryGroup {
	name := beans.Name(beans.NewBeansID().String())

	id, err := u.adapter.CategoryGroupCreate(u.t, u.context, name)
	require.NoError(u.t, err)
	group, err := u.adapter.CategoryGroupGet(u.t, u.context, id)
	require.NoError(u.t, err)

	return group
}

func (u *userAndBudget) Category(opt specification.CategoryOpts) beans.Category {
	name := beans.Name(beans.NewBeansID().String())

	id, err := u.adapter.CategoryCreate(u.t, u.context, opt.Group.ID, name)
	require.NoError(u.t, err)
	category, err := u.adapter.CategoryGet(u.t, u.context, id)
	require.NoError(u.t, err)

	return category
}

func (u *userAndBudget) Transaction(opt specification.TransactionOpts) beans.Transaction {
	if opt.Date.Empty() {
		opt.Date = beans.NewDate(testutils.RandomTime())
	}
	id, err := u.adapter.TransactionCreate(u.t, u.context, beans.TransactionCreateParams{
		TransactionParams: beans.TransactionParams{
			AccountID:  opt.Account.ID,
			CategoryID: opt.Category.ID,
			Amount:     opt.Amount,
			Date:       opt.Date,
		},
	})
	require.NoError(u.t, err)
	transaction, err := u.adapter.TransactionGet(u.t, u.context, id)
	require.NoError(u.t, err)

	return transaction.Transaction
}

// Test

func (a *httpAdapter) UserAndBudget(t *testing.T) specification.TestUserAndBudget {
	ctx := specification.Context{}

	// make new user
	username := beans.NewBeansID().String()
	err := a.UserRegister(t, ctx, beans.Username(username), beans.Password("password"))
	require.Nil(t, err)

	// login as user
	sessionID, err := a.UserLogin(t, ctx, beans.Username(username), beans.Password("password"))
	require.Nil(t, err)

	ctx.SessionID = sessionID

	// make budget
	budgetID, err := a.BudgetCreate(t, ctx, beans.Name(beans.NewBeansID().String()))
	require.Nil(t, err)
	budget, err := a.BudgetGet(t, ctx, budgetID)
	require.Nil(t, err)

	ctx.BudgetID = budget.ID

	return &userAndBudget{
		t:         t,
		sessionID: sessionID,
		budget:    budget,
		context:   ctx,
		adapter:   a,
	}
}
