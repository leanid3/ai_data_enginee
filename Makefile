# Makefile для проекта "Инженер данных"

.PHONY: help build up down logs clean test proto

# Цвета для вывода
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

help: ## Показать справку
	@echo "$(GREEN)Доступные команды:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Собрать все Docker образы
	@echo "$(GREEN)Сборка всех сервисов...$(NC)"
	docker-compose build

up: ## Запустить все сервисы
	@echo "$(GREEN)Запуск всех сервисов...$(NC)"
	docker-compose up -d
	@echo "$(GREEN)Сервисы запущены!$(NC)"
	@echo "$(YELLOW)Доступные URL:$(NC)"
	@echo "  API Gateway: http://localhost:8080"
	@echo "  File Service: http://localhost:8081"
	@echo "  Airflow: http://localhost:8082 (admin/admin)"
	@echo "  MinIO: http://localhost:9001 (minioadmin/minioadmin)"
	@echo "  Ollama: http://localhost:11434"

down: ## Остановить все сервисы
	@echo "$(GREEN)Остановка всех сервисов...$(NC)"
	docker-compose down

restart: down up ## Перезапустить все сервисы

logs: ## Показать логи всех сервисов
	docker-compose logs -f

logs-api: ## Показать логи API Gateway
	docker-compose logs -f api-gateway

logs-file: ## Показать логи File Service
	docker-compose logs -f file-service

logs-chat: ## Показать логи Chat Service
	docker-compose logs -f chat-service

logs-llm: ## Показать логи LLM Service
	docker-compose logs -f llm-service

logs-orchestration: ## Показать логи Orchestration Service
	docker-compose logs -f orchestration-service

logs-airflow: ## Показать логи Airflow
	docker-compose logs -f airflow-webserver airflow-scheduler

status: ## Показать статус всех сервисов
	@echo "$(GREEN)Статус сервисов:$(NC)"
	docker-compose ps

clean: ## Очистить все контейнеры и volumes
	@echo "$(RED)Очистка всех контейнеров и volumes...$(NC)"
	docker-compose down -v --remove-orphans
	docker system prune -f

proto: ## Генерировать protobuf файлы для всех сервисов
	@echo "$(GREEN)Генерация protobuf файлов...$(NC)"
	@for service in api-geteway file-service chat-service llm-service orchestration-service; do \
		if [ -d "$$service" ]; then \
			echo "Генерация для $$service..."; \
			cd $$service && \
			protoc --go_out=. --go_opt=paths=source_relative \
				--go-grpc_out=. --go-grpc_opt=paths=source_relative \
				--grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
				proto/*.proto && \
			cd ..; \
		fi; \
	done

test: ## Запустить тесты
	@echo "$(GREEN)Запуск тестов...$(NC)"
	@for service in api-geteway file-service chat-service llm-service orchestration-service; do \
		if [ -d "$$service" ]; then \
			echo "Тестирование $$service..."; \
			cd $$service && go test ./... && cd ..; \
		fi; \
	done

test-quick: ## Быстрое тестирование системы
	@echo "$(GREEN)Быстрое тестирование системы...$(NC)"
	./quick-test.sh

test-workflow: ## Полное тестирование workflow
	@echo "$(GREEN)Полное тестирование workflow...$(NC)"
	./test-workflow.sh

test-system: ## Проверка состояния системы
	@echo "$(GREEN)Проверка состояния системы...$(NC)"
	./check-system.sh

demo: ## Демонстрация системы
	@echo "$(GREEN)Демонстрация системы...$(NC)"
	./demo.sh

test-api: ## Тестирование API Gateway
	@echo "$(GREEN)Тестирование API Gateway...$(NC)"
	curl -X GET http://localhost:8080/health || echo "API Gateway недоступен"

test-file: ## Тестирование File Service
	@echo "$(GREEN)Тестирование File Service...$(NC)"
	curl -X POST http://localhost:8081/api/v1/files/upload \
		-H "Content-Type: multipart/form-data" \
		-F "file=@test.csv" \
		-F "user_id=test-user" || echo "File Service недоступен"

test-airflow: ## Тестирование Airflow
	@echo "$(GREEN)Тестирование Airflow...$(NC)"
	curl -X GET http://localhost:8082/health || echo "Airflow недоступен"

install-tools: ## Установить необходимые инструменты
	@echo "$(GREEN)Установка инструментов...$(NC)"
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

dev: ## Запуск в режиме разработки
	@echo "$(GREEN)Запуск в режиме разработки...$(NC)"
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

# Специальные команды для разработки
dev-logs: ## Логи в режиме разработки
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml logs -f

dev-down: ## Остановка в режиме разработки
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml down

# Команды для мониторинга
monitor: ## Мониторинг ресурсов
	@echo "$(GREEN)Мониторинг ресурсов:$(NC)"
	docker stats

# Команды для базы данных
db-shell: ## Подключение к PostgreSQL
	docker-compose exec postgres psql -U testuser -d testdb

db-backup: ## Резервное копирование БД
	@echo "$(GREEN)Создание резервной копии БД...$(NC)"
	docker-compose exec postgres pg_dump -U testuser testdb > backup_$(shell date +%Y%m%d_%H%M%S).sql

# Команды для MinIO
minio-shell: ## Подключение к MinIO
	docker-compose exec minio mc

# Команды для Airflow
airflow-dags: ## Показать DAGs в Airflow
	@echo "$(GREEN)DAGs в Airflow:$(NC)"
	curl -s http://localhost:8082/api/v1/dags | jq '.dags[].dag_id'

airflow-tasks: ## Показать задачи в Airflow
	@echo "$(GREEN)Задачи в Airflow:$(NC)"
	curl -s http://localhost:8082/api/v1/dags/data_analysis_pipeline/tasks | jq '.tasks[].task_id'

# Команды для LLM
llm-models: ## Показать доступные модели LLM
	@echo "$(GREEN)Доступные модели LLM:$(NC)"
	curl -s http://localhost:11434/api/tags | jq '.models[].name'

# Команды для очистки
clean-logs: ## Очистить логи
	@echo "$(GREEN)Очистка логов...$(NC)"
	docker-compose logs --tail=0 -f > /dev/null

clean-volumes: ## Очистить volumes
	@echo "$(RED)Очистка volumes...$(NC)"
	docker-compose down -v

# Команды для обновления
update: ## Обновить все образы
	@echo "$(GREEN)Обновление образов...$(NC)"
	docker-compose pull
	docker-compose up -d

# Команды для масштабирования
scale-api: ## Масштабировать API Gateway
	docker-compose up -d --scale api-gateway=2

scale-file: ## Масштабировать File Service
	docker-compose up -d --scale file-service=2

# Команды для безопасности
security-scan: ## Сканирование безопасности
	@echo "$(GREEN)Сканирование безопасности...$(NC)"
	docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
		aquasec/trivy image $(shell docker-compose config | grep image | head -1 | awk '{print $$2}')

# Команды для документации
docs: ## Генерация документации
	@echo "$(GREEN)Генерация документации...$(NC)"
	@if command -v swagger >/dev/null 2>&1; then \
		swagger generate spec -o swagger.json; \
		echo "Документация сгенерирована: swagger.json"; \
	else \
		echo "Swagger не установлен. Установите: go install github.com/go-swagger/go-swagger/cmd/swagger@latest"; \
	fi

# Команды для развертывания
deploy: build up ## Полное развертывание
	@echo "$(GREEN)Развертывание завершено!$(NC)"

# Команды для отладки
debug: ## Отладка сервисов
	@echo "$(GREEN)Информация для отладки:$(NC)"
	@echo "Docker версия: $(shell docker --version)"
	@echo "Docker Compose версия: $(shell docker-compose --version)"
	@echo "Go версия: $(shell go version)"
	@echo "Порты: $(shell netstat -tulpn | grep -E ':(8080|8081|8082|9000|9001|11434)' | head -10)"