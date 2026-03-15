include .env
export $(shell sed 's/=.*//' .env)

up:
	docker compose up --build

down:
	docker compose down

dev-start:
	docker run --name fyp-postgres \
		-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
		-e POSTGRES_USER=$(DB_USER) \
		-e POSTGRES_DB=$(DB_NAME) \
		-p $(DB_PORT):5432 \
		-d postgres:15

dev:
	go run -race backend/main.go server

mcp-one:
	python mcp_server_1/server.py

