#!/bin/bash

# Скрипт для проверки состояния системы "Инженер данных"

set -e

# Цвета для вывода
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${BLUE}"
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║              ПРОВЕРКА СИСТЕМЫ 'ИНЖЕНЕР ДАННЫХ'               ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

# Проверка Docker
check_docker() {
    print_info "Проверка Docker..."
    if command -v docker &> /dev/null; then
        print_success "Docker установлен: $(docker --version)"
    else
        print_error "Docker не установлен"
        exit 1
    fi
    
    if command -v docker-compose &> /dev/null; then
        print_success "Docker Compose установлен: $(docker-compose --version)"
    else
        print_error "Docker Compose не установлен"
        exit 1
    fi
}

# Проверка контейнеров
check_containers() {
    print_info "Проверка контейнеров..."
    
    # Список ожидаемых сервисов
    services=("api-gateway" "file-service" "chat-service" "llm-service" "orchestration-service" "postgres" "minio" "redis" "ollama" "airflow-webserver" "airflow-scheduler")
    
    for service in "${services[@]}"; do
        if docker-compose ps | grep -q "$service.*Up"; then
            print_success "$service: Запущен"
        else
            print_warning "$service: Не запущен"
        fi
    done
}

# Проверка портов
check_ports() {
    print_info "Проверка портов..."
    
    ports=("8080" "8081" "8082" "50051" "50054" "50055" "50056" "50057" "5432" "6379" "9000" "9001" "11434")
    
    for port in "${ports[@]}"; do
        if netstat -tulpn 2>/dev/null | grep -q ":$port "; then
            print_success "Порт $port: Открыт"
        else
            print_warning "Порт $port: Закрыт"
        fi
    done
}

# Проверка HTTP эндпоинтов
check_endpoints() {
    print_info "Проверка HTTP эндпоинтов..."
    
    # API Gateway
    if curl -s http://localhost:8080/health >/dev/null 2>&1; then
        print_success "API Gateway: Доступен"
    else
        print_warning "API Gateway: Недоступен"
    fi
    
    # File Service
    if curl -s http://localhost:8081/health >/dev/null 2>&1; then
        print_success "File Service: Доступен"
    else
        print_warning "File Service: Недоступен"
    fi
    
    # Airflow
    if curl -s http://localhost:8082/health >/dev/null 2>&1; then
        print_success "Airflow: Доступен"
    else
        print_warning "Airflow: Недоступен"
    fi
    
    # MinIO
    if curl -s http://localhost:9000/minio/health/live >/dev/null 2>&1; then
        print_success "MinIO: Доступен"
    else
        print_warning "MinIO: Недоступен"
    fi
    
    # Ollama
    if curl -s http://localhost:11434/api/tags >/dev/null 2>&1; then
        print_success "Ollama: Доступен"
    else
        print_warning "Ollama: Недоступен"
    fi
}

# Проверка логов на ошибки
check_logs() {
    print_info "Проверка логов на ошибки..."
    
    # Проверяем логи каждого сервиса
    services=("api-gateway" "file-service" "chat-service" "llm-service" "orchestration-service")
    
    for service in "${services[@]}"; do
        if docker-compose logs --tail=10 "$service" 2>/dev/null | grep -i error >/dev/null; then
            print_warning "$service: Найдены ошибки в логах"
        else
            print_success "$service: Логи чистые"
        fi
    done
}

# Проверка ресурсов
check_resources() {
    print_info "Проверка ресурсов..."
    
    # Проверяем использование памяти
    memory_usage=$(docker stats --no-stream --format "table {{.MemUsage}}" | tail -n +2 | awk '{sum += $1} END {print sum}')
    if [ -n "$memory_usage" ]; then
        print_info "Использование памяти: $memory_usage"
    fi
    
    # Проверяем использование CPU
    cpu_usage=$(docker stats --no-stream --format "table {{.CPUPerc}}" | tail -n +2 | awk '{sum += $1} END {print sum}')
    if [ -n "$cpu_usage" ]; then
        print_info "Использование CPU: $cpu_usage"
    fi
}

# Проверка базы данных
check_database() {
    print_info "Проверка базы данных..."
    
    if docker-compose exec -T postgres pg_isready -U testuser >/dev/null 2>&1; then
        print_success "PostgreSQL: Доступна"
    else
        print_warning "PostgreSQL: Недоступна"
    fi
}

# Проверка хранилища
check_storage() {
    print_info "Проверка хранилища..."
    
    if docker-compose exec -T minio mc ls minio/files >/dev/null 2>&1; then
        print_success "MinIO: Доступен"
    else
        print_warning "MinIO: Недоступен"
    fi
}

# Основная функция
main() {
    print_header
    
    check_docker
    echo ""
    
    check_containers
    echo ""
    
    check_ports
    echo ""
    
    check_endpoints
    echo ""
    
    check_logs
    echo ""
    
    check_database
    echo ""
    
    check_storage
    echo ""
    
    check_resources
    echo ""
    
    print_info "Проверка завершена!"
    echo ""
    echo "Для получения подробной информации используйте:"
    echo "  docker-compose logs -f [service-name]"
    echo "  docker-compose ps"
    echo "  docker stats"
}

# Запуск основной функции
main "$@"