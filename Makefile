# Variables
GO_BIN = ~/go/bin
APP_NAME = piggytrack
SRC_DIR = ./cmd/piggytrack
AIR = ~/go/bin/air
GOOSE = goose
SQLC = sqlc
MIGRATIONS_DIR = ./internal/db/migrations
WEB_DIR = ./web

include .env
export

.PHONY: dev
dev:
	${AIR} serve

.PHONY: build
build:
	go build -o $(APP_NAME) $(SRC_DIR)

.PHONY: migrate-up
migrate-up:
	$(GOOSE) -dir=${MIGRATIONS_DIR} postgres $(DATABASE_URL) up

.PHONY: migrate-down
migrate-down:
	$(GOOSE) -dir=${MIGRATIONS_DIR} postgres $(DATABASE_URL) down

.PHONY: migrate-reset
migrate-reset:
	$(GOOSE) -dir=${MIGRATIONS_DIR} postgres $(DATABASE_URL) reset

.PHONY: migrate-status
migrate-status:
	$(GOOSE) -dir=${MIGRATIONS_DIR} postgres $(DATABASE_URL) status
