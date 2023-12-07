.PHONY: default build run migrate migration gensql test

default: run

SERVER_DIR = ./server/

run:
	cd $(SERVER_DIR) && go run ./cmd/beansd

migrate:
	cd $(SERVER_DIR) && go run ./cmd/beans migrate

migration:
	cd $(SERVER_DIR) && migrate create -dir internal/sql/migrations -ext sql ${NAME}
	cd $(SERVER_DIR) && rm internal/sql/migrations/*.down.sql

gensql:
	cd $(SERVER_DIR) && sqlc generate

genmock:
	cd $(SERVER_DIR) && go generate ./...

test:
	cd $(SERVER_DIR) && go test -tags test --count=1 ./... 


