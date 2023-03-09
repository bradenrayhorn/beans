default: build

build:
	@go build -o ./beansd ./cmd/beansd

run:
	@go run ./cmd/beansd

migrate:
	@go run ./cmd/beans migrate

migration:
	@migrate create -dir internal/sql/migrations -ext sql ${NAME}
	@rm internal/sql/migrations/*.down.sql

gensql:
	@sqlc generate

genmock:
	@go generate ./...

test:
	@go test -tags test --count=1 ./... 

.PHONY: default build run migrate migration gensql test

