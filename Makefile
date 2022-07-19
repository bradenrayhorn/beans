default: build

build:
	@go build -o ./beansd ./cmd/beansd

run:
	@go run ./cmd/beansd

migrate:
	@migrate -database "postgres://postgres:password@127.0.0.1:5432/beans?sslmode=disable" -path internal/sql/migrations up

migration:
	@migrate create -dir internal/sql/migrations -ext sql ${NAME}
	@rm internal/sql/migrations/*.down.sql

gensql:
	@sqlc generate

test:
	@go test -v --count=1 ./... 

.PHONY: default build run migrate migration gensql test

