package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/bradenrayhorn/beans/cmd/beansd"
	"github.com/orlangure/gnomock"
	pg "github.com/orlangure/gnomock/preset/postgres"
)

type TestApplication struct {
	application       *main.Application
	postgresContainer *gnomock.Container
}

func StartApplication(tb testing.TB) *TestApplication {
	p := pg.Preset(
		pg.WithVersion("13.4"),
		pg.WithDatabase("beans"),
		pg.WithQueriesFile("../../internal/sql/migrations/20220708015254_create_users_table.up.sql"),
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

func (ta *TestApplication) PostRequest(path string, body map[string]interface{}) (*http.Response, error) {
	jsonBody, _ := json.Marshal(body)
	boundAddr := ta.application.HttpServer().GetBoundAddr()
	return http.Post(fmt.Sprintf("http://%s/%s", boundAddr, path), "application/json", bytes.NewReader(jsonBody))
}
