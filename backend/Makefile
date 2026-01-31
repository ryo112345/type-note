GOLANGCI_LINT ?= $(shell go env GOPATH)/bin/golangci-lint
GO_TEST ?= go test
TEST_FLAGS ?= -count=1

MIGRATE ?= migrate
MIGRATIONS_DIR ?= migrations
NAME ?=
DB_HOST ?= localhost
ENV_FILE ?= .env

ifneq (,$(wildcard .env))
include .env
export $(shell grep -E '^[A-Za-z0-9_]+=' .env | cut -d= -f1)
endif

# Build DATABASE_URL after .env is loaded
DATABASE_URL ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

.PHONY: lint
lint:
	@$(GOLANGCI_LINT) run ./...

.PHONY: test
test:
	@$(GO_TEST) ./...

.PHONY: test-all
test-all:
	@$(GO_TEST) $(TEST_FLAGS) ./...

.PHONY: test-pkg
test-pkg:
	@test -n "$(PKG)" || (echo "PKG is required, e.g. make test-pkg PKG=./internal/domain/note"; exit 1)
	@$(GO_TEST) $(TEST_FLAGS) $(PKG)

.PHONY: build
build:
	@go build ./...

.PHONY: build-api
build-api:
	@go build -o bin/api ./cmd/api

.PHONY: run-api
run-api:
	@go run ./cmd/api/main.go

# Database migrations
.PHONY: migrate-up
migrate-up:
	@bash -c 'set -a; if [ -f $(ENV_FILE) ]; then source $(ENV_FILE); fi; set +a; \
	if [ -z "$$MIGRATION_DATABASE_URL" ]; then echo "Set MIGRATION_DATABASE_URL in $(ENV_FILE) or environment"; exit 1; fi; \
	"$(MIGRATE)" -path "$(MIGRATIONS_DIR)" -database "$$MIGRATION_DATABASE_URL" up'

.PHONY: migrate-down
migrate-down:
	@bash -c 'set -a; if [ -f $(ENV_FILE) ]; then source $(ENV_FILE); fi; set +a; \
	if [ -z "$$MIGRATION_DATABASE_URL" ]; then echo "Set MIGRATION_DATABASE_URL in $(ENV_FILE) or environment"; exit 1; fi; \
	"$(MIGRATE)" -path "$(MIGRATIONS_DIR)" -database "$$MIGRATION_DATABASE_URL" down 1'

.PHONY: migrate-new
migrate-new:
	@test -n "$(NAME)" || (echo "NAME is required, e.g. make migrate-new NAME=add_users"; exit 1)
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) $(NAME)

# sqlc code generation
.PHONY: sqlc
sqlc:
	@sqlc generate
