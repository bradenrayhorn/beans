package main

import (
	"errors"
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	PostgresURL string

	Port string
}

var k = koanf.New(".")

func LoadConfig() (Config, error) {

	// load defaults
	err := k.Load(confmap.Provider(map[string]interface{}{
		"http.port": "8000",
		"pg.url":    "postgres://postgres:password@127.0.0.1:5432/beans",
	}, "."), nil)
	if err != nil {
		return Config{}, err
	}

	// load from dotenv
	if err := k.Load(file.Provider(".env"), dotenv.Parser()); err != nil {
		if !errors.Is(err, os.ErrNotExist) { // ignore .env not found error
			return Config{}, err
		}
	}
	if err != nil {
		return Config{}, err
	}

	// load from env
	err = k.Load(env.Provider("BEANS_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "BEANS_")), "_", ".", -1)
	}), nil)
	if err != nil {
		return Config{}, err
	}

	return Config{
		PostgresURL: k.String("pg.url"),
		Port:        k.String("http.port"),
	}, nil
}
