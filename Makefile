default: build

build:
	@go build -o ./beansd ./cmd/beansd

.PHONY: default build

