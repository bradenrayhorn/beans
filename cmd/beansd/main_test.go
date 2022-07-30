package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/cmd/beansd"
	"github.com/orlangure/gnomock"
	pg "github.com/orlangure/gnomock/preset/postgres"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type TestApplication struct {
	application       *main.Application
	postgresContainer *gnomock.Container
}

func StartApplication(tb testing.TB) *TestApplication {
	p := pg.Preset(
		pg.WithVersion("13.4"),
		pg.WithDatabase("beans"),
		pg.WithQueries(getMigrationQueries(tb)),
	)

	container, err := gnomock.Start(p)
	if err != nil {
		tb.Fatal(err)
	}

	testApp := &TestApplication{postgresContainer: container}
	testApp.application = main.NewApplication(main.Config{
		Postgres: main.PostgresConfig{
			Addr:     fmt.Sprintf("%s:%d", "127.0.0.1", container.DefaultPort()),
			Username: "postgres",
			Password: "password",
			Database: "beans",
		},
		Port: "0",
	})

	if err := testApp.application.Start(); err != nil {
		tb.Fatal(err)
	}

	return testApp
}

func (ta *TestApplication) Stop(tb testing.TB) {
	if err := ta.application.Stop(); err != nil {
		tb.Fatal(err)
	}

	gnomock.Stop(ta.postgresContainer)
}

func getMigrationQueries(tb testing.TB) string {
	queries := ""
	err := filepath.WalkDir("../../internal/sql/migrations/", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		queries += string(content)

		return nil
	})
	if err != nil {
		tb.Fatal(err)
	}

	return queries
}

// http request helpers

func (ta *TestApplication) PostRequest(tb testing.TB, path string, options *RequestOptions) *TestResponse {
	return ta.doRequest(tb, "POST", path, options)
}

func (ta *TestApplication) GetRequest(tb testing.TB, path string, options *RequestOptions) *TestResponse {
	return ta.doRequest(tb, "GET", path, options)
}

type RequestOptions struct {
	SessionID string
	Body      any
}

type TestResponse struct {
	resp                *http.Response
	StatusCode          int
	Body                string
	SessionIDFromCookie string
}

func (ta *TestApplication) doRequest(tb testing.TB, method string, path string, options *RequestOptions) *TestResponse {
	boundAddr := ta.application.HttpServer().GetBoundAddr()
	url := fmt.Sprintf("http://%s/%s", boundAddr, path)

	if options == nil {
		options = &RequestOptions{}
	}

	var body io.Reader = nil
	switch options.Body.(type) {
	case string:
		body = bytes.NewReader([]byte(options.Body.(string)))
	case nil:
		body = nil
	default:
		reqBytes, _ := json.Marshal(options.Body)
		body = bytes.NewReader(reqBytes)
	}
	request, err := http.NewRequest(method, url, body)
	require.Nil(tb, err)

	if len(options.SessionID) > 0 {
		request.AddCookie(&http.Cookie{Name: "session_id", Value: options.SessionID})
	}

	client := http.Client{}
	resp, err := client.Do(request)
	require.Nil(tb, err)

	respBody, err := io.ReadAll(resp.Body)
	require.Nil(tb, err)

	var sessionID string
	for _, c := range resp.Cookies() {
		if c.Valid() == nil && c.Name == "session_id" {
			sessionID = c.Value
		}
	}

	return &TestResponse{resp: resp, StatusCode: resp.StatusCode, Body: string(respBody), SessionIDFromCookie: sessionID}
}

// database helpers

func (ta *TestApplication) CreateUser(tb testing.TB, username string, password string) *beans.User {
	userID := beans.UserID(ksuid.New())
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.Nil(tb, err)
	err = ta.application.UserRepository().Create(context.Background(), userID, beans.Username(username), beans.PasswordHash(passwordHash))
	require.Nil(tb, err)
	return &beans.User{ID: userID, Username: beans.Username(username), PasswordHash: beans.PasswordHash(passwordHash)}
}

func (ta *TestApplication) CreateSession(tb testing.TB, user *beans.User) *beans.Session {
	session, err := ta.application.SessionRepository().Create(user.ID)
	require.Nil(tb, err)
	return session
}

func (ta *TestApplication) CreateUserAndSession(tb testing.TB) (*beans.User, *beans.Session) {
	user := ta.CreateUser(tb, "testuser", "password")
	session := ta.CreateSession(tb, user)
	return user, session
}

func (ta *TestApplication) CreateBudget(tb testing.TB, name string, user *beans.User) *beans.Budget {
	id := beans.NewBeansID()
	err := ta.application.BudgetRepository().Create(context.Background(), id, beans.BudgetName(name), user.ID)
	require.Nil(tb, err)
	return &beans.Budget{ID: id, Name: beans.BudgetName(name)}
}
