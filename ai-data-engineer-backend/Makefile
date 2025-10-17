.PHONY: help build run test docker-up docker-down clean lint fmt deps

# Переменные
BINARY_NAME=ai-data-engineer-backend
DOCKER_IMAGE=ai-data-engineer-backend
DOCKER_TAG=latest
CONFIG_PATH=configs/config.yaml

# Цвета для вывода
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

help: ## Показать help
	@echo "$(BLUE)AI Data Engineer Backend - Available Commands$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

build: ## Собрать Go binary
	@echo "$(BLUE)Building $(BINARY_NAME)...$(NC)"
	@go build -o $(BINARY_NAME) ./cmd/server
	@echo "$(GREEN)Build completed!$(NC)"

run: ## Запустить приложение локально
	@echo "$(BLUE)Starting $(BINARY_NAME)...$(NC)"
	@CONFIG_PATH=$(CONFIG_PATH) ./$(BINARY_NAME)

test: ## Запустить unit тесты
	@echo "$(BLUE)Running tests...$(NC)"
	@go test -v ./...

test-integration: ## Запустить integration тесты
	@echo "$(BLUE)Running integration tests...$(NC)"
	@go test -v -tags=integration ./...

docker-up: ## Поднять все сервисы в Docker
	@echo "$(BLUE)Starting Docker services...$(NC)"
	@docker-compose up -d
	@echo "$(GREEN)Docker services started!$(NC)"

docker-down: ## Остановить все Docker сервисы
	@echo "$(BLUE)Stopping Docker services...$(NC)"
	@docker-compose down
	@echo "$(GREEN)Docker services stopped!$(NC)"

docker-logs: ## Показать логи всех сервисов
	@echo "$(BLUE)Showing Docker logs...$(NC)"
	@docker-compose logs -f

clean: ## Очистить build артефакты
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	@rm -f $(BINARY_NAME)
	@go clean
	@echo "$(GREEN)Clean completed!$(NC)"

lint: ## Запустить golangci-lint
	@echo "$(BLUE)Running linter...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)golangci-lint not found. Installing...$(NC)"; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.54.2; \
		golangci-lint run; \
	fi

fmt: ## Форматировать код с gofmt
	@echo "$(BLUE)Formatting code...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)Code formatted!$(NC)"

deps: ## Установить зависимости
	@echo "$(BLUE)Installing dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)Dependencies installed!$(NC)"

dev: ## Запустить в режиме разработки с hot reload
	@echo "$(BLUE)Starting development server...$(NC)"
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "$(YELLOW)Air not found. Installing...$(NC)"; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

docker-build: ## Собрать Docker образ
	@echo "$(BLUE)Building Docker image...$(NC)"
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "$(GREEN)Docker image built!$(NC)"

docker-run: ## Запустить Docker контейнер
	@echo "$(BLUE)Running Docker container...$(NC)"
	@docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

check-tools: ## Проверить наличие необходимых инструментов
	@echo "$(BLUE)Checking required tools...$(NC)"
	@command -v go >/dev/null 2>&1 || { echo "$(RED)Go is not installed$(NC)"; exit 1; }
	@command -v docker >/dev/null 2>&1 || { echo "$(RED)Docker is not installed$(NC)"; exit 1; }
	@command -v docker-compose >/dev/null 2>&1 || { echo "$(RED)Docker Compose is not installed$(NC)"; exit 1; }
	@echo "$(GREEN)All required tools are available!$(NC)"

setup: check-tools deps ## Настроить проект
	@echo "$(GREEN)Project setup completed!$(NC)"

# Команда по умолчанию
.DEFAULT_GOAL := help

