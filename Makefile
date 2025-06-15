GO_VERSION := 1.24
PROJECT_NAME := testing

# Colors for output
GREEN := \033[32m
BLUE := \033[34m
YELLOW := \033[33m
RED := \033[31m
RESET := \033[0m

.PHONY: help
help: ## Show this help message
	@echo "$(BLUE)Iron CLI Development Commands$(RESET)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)%-15s$(RESET) %s\n", $$1, $$2}'

.PHONY: clean
clean: ## Remove the testing directory
	@echo "$(YELLOW)🧹 Cleaning up testing directory...$(RESET)"
	@if [ -d "$(PROJECT_NAME)" ]; then \
		rm -rf "$(PROJECT_NAME)" && \
		echo "$(GREEN)✓ Removed existing '$(PROJECT_NAME)' directory$(RESET)"; \
	else \
		echo "$(BLUE)ℹ Directory '$(PROJECT_NAME)' doesn't exist$(RESET)"; \
	fi

.PHONY: generate
generate: clean ## Generate OAuth project, clean up first
	@echo "$(YELLOW)🔧 Generating OAuth project...$(RESET)"
	@go run main.go generate oauth $(PROJECT_NAME)
	@echo "$(GREEN)✓ OAuth project generated successfully$(RESET)"

.PHONY: dev
dev: generate ## Generate project and open in VS Code
	@echo "$(YELLOW)📂 Opening project in VS Code...$(RESET)"
	@cd $(PROJECT_NAME) && code .
	@echo "$(GREEN)✓ Project opened in VS Code$(RESET)"

.PHONY: generate-custom
generate-custom: ## Generate OAuth project with custom name (usage: make generate-custom PROJECT_NAME=myproject)
	@echo "$(YELLOW)🧹 Cleaning up $(PROJECT_NAME) directory...$(RESET)"
	@if [ -d "$(PROJECT_NAME)" ]; then \
		rm -rf "$(PROJECT_NAME)" && \
		echo "$(GREEN)✓ Removed existing '$(PROJECT_NAME)' directory$(RESET)"; \
	fi
	@echo "$(YELLOW)🔧 Generating OAuth project '$(PROJECT_NAME)'...$(RESET)"
	@go run main.go generate oauth $(PROJECT_NAME)
	@echo "$(YELLOW)📂 Opening project in VS Code...$(RESET)"
	@cd $(PROJECT_NAME) && code .
	@echo "$(GREEN)✓ Project '$(PROJECT_NAME)' generated and opened$(RESET)"

.PHONY: build
build: ## Build the iron binary
	@echo "$(YELLOW)🔨 Building iron binary...$(RESET)"
	@go build -o bin/iron main.go
	@echo "$(GREEN)✓ Binary built at bin/iron$(RESET)"

.PHONY: install
install: ## Install iron globally
	@echo "$(YELLOW)📦 Installing iron globally...$(RESET)"
	@go install .
	@echo "$(GREEN)✓ Iron installed globally$(RESET)"

.PHONY: lint
lint: ## Run golangci-lint
	@echo "$(YELLOW)🔍 Running linter...$(RESET)"
	@golangci-lint run
	@echo "$(GREEN)✓ Linting completed$(RESET)"

.PHONY: test
test: ## Run tests
	@echo "$(YELLOW)🧪 Running tests...$(RESET)"
	@go test ./...
	@echo "$(GREEN)✓ Tests completed$(RESET)"

.PHONY: fmt
fmt: ## Format code
	@echo "$(YELLOW)✨ Formatting code...$(RESET)"
	@go fmt ./...
	@echo "$(GREEN)✓ Code formatted$(RESET)"

.PHONY: mod-tidy
mod-tidy: ## Clean up go.mod
	@echo "$(YELLOW)📝 Tidying go.mod...$(RESET)"
	@go mod tidy
	@echo "$(GREEN)✓ go.mod tidied$(RESET)"

.PHONY: all
all: fmt lint test build ## Run fmt, lint, test, and build

# Default target
.DEFAULT_GOAL := help