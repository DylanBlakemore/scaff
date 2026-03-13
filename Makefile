.PHONY: help format lint test security build ci

GO ?= go
BIN_DIR ?= bin
BINARY ?= scaff
GOOS := $(shell $(GO) env GOOS)
EXE :=
ifeq ($(GOOS),windows)
EXE = .exe
endif

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-12s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

format: ## Format code
	$(GO) fmt ./...

lint: ## Lint (style/static analysis)
	$(GO) run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run

test: ## Run tests
	$(GO) test ./...

security: ## Security + dependency checks
	$(GO) vet ./...
	$(GO) mod verify
	$(GO) run golang.org/x/vuln/cmd/govulncheck@latest ./...

build: ## Build local binary to ./bin/
	@mkdir -p "$(BIN_DIR)"
	CGO_ENABLED=0 $(GO) build -o "$(BIN_DIR)/$(BINARY)$(EXE)" .

ci: format lint test security ## Run everything (CI parity)
