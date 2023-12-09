package http_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	goHttp "net/http"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/http"
	"github.com/bradenrayhorn/beans/server/internal/mocks"
)

type httpTest struct {
	httpServer *http.Server

	accountContract     *mocks.MockAccountContract
	budgetContract      *mocks.MockBudgetContract
	categoryContract    *mocks.MockCategoryContract
	monthContract       *mocks.MockMonthContract
	payeeContract       *mocks.MockPayeeContract
	transactionContract *mocks.MockTransactionContract
	userContract        *mocks.MockUserContract
}

func newHttpTest(tb testing.TB) *httpTest {
	test := &httpTest{
		accountContract:     mocks.NewMockAccountContract(),
		budgetContract:      mocks.NewMockBudgetContract(),
		categoryContract:    mocks.NewMockCategoryContract(),
		monthContract:       mocks.NewMockMonthContract(),
		payeeContract:       mocks.NewMockPayeeContract(),
		transactionContract: mocks.NewMockTransactionContract(),
		userContract:        mocks.NewMockUserContract(),
	}

	test.httpServer = http.NewServer(
		test.accountContract,
		test.budgetContract,
		test.categoryContract,
		test.monthContract,
		test.payeeContract,
		test.transactionContract,
		test.userContract,
	)

	if err := test.httpServer.Open(":0"); err != nil {
		tb.Fatal(err)
	}

	return test
}

func (t *httpTest) Stop(tb testing.TB) {
	if err := t.httpServer.Close(); err != nil {
		tb.Fatal(err)
	}
}

type HTTPRequest struct {
	method string
	path   string
	body   any
	user   *beans.User
	budget *beans.Budget
}

type HTTPResponse struct {
	*goHttp.Response
	body string
}

func (t *httpTest) DoRequest(tb testing.TB, httpRequest HTTPRequest) *HTTPResponse {
	var reqBody io.Reader
	switch body := httpRequest.body.(type) {
	case string:
		reqBody = bytes.NewReader([]byte(body))
	default:
		reqBody = nil
	}

	req, err := goHttp.NewRequest(
		httpRequest.method,
		"http://"+t.httpServer.GetBoundAddr()+httpRequest.path,
		reqBody,
	)
	if err != nil {
		tb.Fatal(err)
	}

	if user := httpRequest.user; user != nil {
		cookie := goHttp.Cookie{
			Name:  "session_id",
			Value: "12345",
		}
		req.AddCookie(&cookie)

		t.userContract.GetAuthFunc.PushHook(func(ctx context.Context, si beans.SessionID) (*beans.AuthContext, error) {
			if string(si) == "12345" {
				return beans.NewAuthContext(user.ID, beans.SessionID(cookie.Value)), nil
			} else {
				return nil, errors.New("received invalid session id in mock")
			}
		})
	}

	if budget := httpRequest.budget; budget != nil {
		req.Header.Add("Budget-ID", budget.ID.String())

		t.budgetContract.GetFunc.PushHook(func(ctx context.Context, ac *beans.AuthContext, i beans.ID) (*beans.Budget, error) {
			if i == budget.ID && ac.UserID() == httpRequest.user.ID {
				return budget, nil
			} else {
				return nil, errors.New("received invalid budget in mock")
			}
		})
	}

	client := &goHttp.Client{}

	resp, err := client.Do(req)
	if err != nil {
		tb.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		tb.Fatal(err)
	}

	return &HTTPResponse{
		Response: resp,
		body:     strings.TrimSpace(string(data)),
	}
}
