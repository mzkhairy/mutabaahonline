.PHONY: run build test tidy migrate-up migrate-down migrate-force

# Load env vars from .env if present
ifneq (,$(wildcard .env))
include .env
endif

run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

test:
	go test ./...

tidy:
	go mod tidy

# Migration commands require golang-migrate CLI installed: https://github.com/golang-migrate/migrate
migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1

migrate-force:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" force ${v}
