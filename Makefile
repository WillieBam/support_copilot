include .env
export $(shell sed 's/=.*//' .env)

.PHONY: up down dev-start dev mcp-one

up:
	podman compose up --build

down:
	podman compose down

dev-start:
	podman run --name support_copilot_postgres \
		-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
		-e POSTGRES_USER=$(DB_USER) \
		-e POSTGRES_DB=$(DB_NAME) \
		-p $(DB_PORT):5432 \
		-d docker.io/library/postgres:15
migrate:
	go run backend/main.go migrate
dev:
	go run -race backend/main.go server

mcp-one:
	python mcp_server_1/server.py

MOCKERY = $(shell pwd)/backend/bin/mockery

$(MOCKERY):
	GOBIN=$(shell pwd)/backend/bin go install github.com/vektra/mockery/v2@v2.53.6

generate: $(MOCKERY)
	$(MOCKERY) --all --dir backend/internal/interfaces --output backend/internal/mocks --outpkg mocks --case camel


test:
	go test -v ./backend/...
