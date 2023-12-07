package main

type Config struct {
	Postgres PostgresConfig

	Port string
}

type PostgresConfig struct {
	Addr     string
	Username string
	Password string
	Database string
}
