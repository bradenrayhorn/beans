.PHONY: default run gensql

default: run

SERVER_DIR = ./server/

run:
	cd $(SERVER_DIR) && go run ./cmd/beansd

gensql:
	cd $(SERVER_DIR) && sqlc generate


