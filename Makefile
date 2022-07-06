default: build

build:
	@go build -o ./beansd ./cmd/beansd

run:
	@go run ./cmd/beansd

.PHONY: default build run

