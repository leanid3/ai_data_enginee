#!/bin/bash

# Скрипт для тестирования интеграции всех сервисов
# AI Data Engineer Backend + LLM + Frontend

set -e

echo "🚀 Запуск тестирования интеграции AI Data Engineer системы"
echo "=================================================="

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Функция для логирования
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Проверка наличия Docker
if ! command -v docker &> /dev/null; then
    error "Docker не установлен!"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    error "Docker Compose не установлен!"
    exit 1
fi

log "Проверка конфигурации docker-compose..."
docker-compose config > /dev/null
success "Конфигурация docker-compose корректна"

# Остановка и удаление существующих контейнеров
log "Остановка существующих контейнеров..."
docker-compose down --remove-orphans

# Сборка образов
log "Сборка Docker образов..."
docker-compose build --no-cache

# Запуск сервисов
log "Запуск сервисов..."
docker-compose up -d

# Ожидание готовности сервисов
log "Ожидание готовности сервисов..."

# Функция для проверки готовности сервиса
wait_for_service() {
    local service_name=$1
    local url=$2
    local max_attempts=30
    local attempt=1
    
    log "Ожидание готовности $service_name..."
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s -f "$url" > /dev/null 2>&1; then
            success "$service_name готов!"
            return 0
        fi
        
        echo -n "."
        sleep 2
        attempt=$((attempt + 1))
    done
    
    error "$service_name не готов после $max_attempts попыток"
    return 1
}

# Проверка готовности сервисов
wait_for_service "PostgreSQL" "http://localhost:5432" || true
wait_for_service "Redis" "http://localhost:6379" || true
wait_for_service "MinIO" "http://localhost:9000/minio/health/live" || true
wait_for_service "ClickHouse" "http://localhost:8123/ping" || true
wait_for_service "LLM Service" "http://localhost:8124/api/v1/health" || true
wait_for_service "AI Backend" "http://localhost:8080/api/v1/health" || true
wait_for_service "Frontend" "http://localhost:3001" || true

# Проверка статуса контейнеров
log "Проверка статуса контейнеров..."
docker-compose ps

# Тестирование API endpoints
log "Тестирование API endpoints..."

# Тест LLM сервиса
log "Тестирование LLM сервиса..."
if curl -s -f "http://localhost:8124/api/v1/health" > /dev/null; then
    success "LLM сервис работает"
else
    warning "LLM сервис недоступен"
fi

# Тест AI Backend
log "Тестирование AI Backend..."
if curl -s -f "http://localhost:8080/api/v1/health" > /dev/null; then
    success "AI Backend работает"
else
    warning "AI Backend недоступен"
fi

# Тест Frontend
log "Тестирование Frontend..."
if curl -s -f "http://localhost:3001" > /dev/null; then
    success "Frontend работает"
else
    warning "Frontend недоступен"
fi

# Проверка логов на ошибки
log "Проверка логов на критические ошибки..."
if docker-compose logs | grep -i "error\|fatal\|panic" | head -5; then
    warning "Обнаружены ошибки в логах"
else
    success "Критических ошибок в логах не найдено"
fi

echo ""
echo "=================================================="
echo "🎉 Тестирование интеграции завершено!"
echo ""
echo "Доступные сервисы:"
echo "  • Frontend: http://localhost:3001"
echo "  • AI Backend API: http://localhost:8080"
echo "  • LLM Service: http://localhost:8124"
echo "  • MinIO Console: http://localhost:9001"
echo "  • ClickHouse: http://localhost:8123"
echo ""
echo "Для просмотра логов: docker-compose logs -f"
echo "Для остановки: docker-compose down"
echo "=================================================="
