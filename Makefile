.PHONY: default run

default: run

SERVER_DIR = ./server/

run:
	cd $(SERVER_DIR) && go run ./cmd/beansd

