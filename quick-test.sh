#!/bin/bash

# Быстрое тестирование отдельных компонентов системы

set -e

# Цвета для вывода
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Переменные
API_BASE="http://localhost:8080"
FILE_SERVICE_BASE="http://localhost:8081"
AIRFLOW_BASE="http://localhost:8082"
TEST_USER="quick-test-$(date +%s)"

print_header() {
    echo -e "${BLUE}"
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║                    БЫСТРОЕ ТЕСТИРОВАНИЕ                      ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

# Функция для тестирования API Gateway
test_api_gateway() {
    print_info "Тестирование API Gateway..."
    
    if curl -s "$API_BASE/health" >/dev/null 2>&1; then
        print_success "API Gateway доступен"
    else
        print_error "API Gateway недоступен"
        return 1
    fi
}

# Функция для тестирования File Service
test_file_service() {
    print_info "Тестирование File Service..."
    
    if curl -s "$FILE_SERVICE_BASE/health" >/dev/null 2>&1; then
        print_success "File Service доступен"
    else
        print_error "File Service недоступен"
        return 1
    fi
}

# Функция для тестирования Airflow
test_airflow() {
    print_info "Тестирование Airflow..."
    
    if curl -s "$AIRFLOW_BASE/health" >/dev/null 2>&1; then
        print_success "Airflow доступен"
    else
        print_error "Airflow недоступен"
        return 1
    fi
}

# Функция для тестирования загрузки небольшого файла
test_small_file_upload() {
    print_info "Тестирование загрузки небольшого файла..."
    
    # Создаем тестовый файл
    echo "id,name,email
1,John Doe,john@example.com
2,Jane Smith,jane@example.com
3,Bob Johnson,bob@example.com" > /tmp/test-small.csv
    
    response=$(curl -s -X POST "$API_BASE/v1/files/upload/csv" \
        -H "Content-Type: multipart/form-data" \
        -F "file=@/tmp/test-small.csv" \
        -F "user_id=$TEST_USER")
    
    if echo "$response" | grep -q "file_id"; then
        FILE_ID=$(echo "$response" | grep -o '"file_id":"[^"]*"' | cut -d'"' -f4)
        print_success "Небольшой файл загружен, ID: $FILE_ID"
        
        # Проверяем информацию о файле
        info_response=$(curl -s -X GET "$API_BASE/v1/files/$FILE_ID?user_id=$TEST_USER")
        if echo "$info_response" | grep -q "file_id"; then
            print_success "Информация о файле получена"
        else
            print_error "Ошибка получения информации о файле"
        fi
    else
        print_error "Ошибка загрузки файла: $response"
        return 1
    fi
    
    # Очищаем тестовый файл
    rm -f /tmp/test-small.csv
}

# Функция для тестирования создания диалога
test_dialog_creation() {
    print_info "Тестирование создания диалога..."
    
    response=$(curl -s -X POST "$API_BASE/v1/dialogs" \
        -H "Content-Type: application/json" \
        -d "{
            \"user_id\": \"$TEST_USER\",
            \"title\": \"Тестовый диалог\",
            \"initial_message\": \"Привет, это тестовое сообщение\"
        }")
    
    if echo "$response" | grep -q "dialog_id"; then
        DIALOG_ID=$(echo "$response" | grep -o '"dialog_id":"[^"]*"' | cut -d'"' -f4)
        print_success "Диалог создан, ID: $DIALOG_ID"
        
        # Проверяем историю диалога
        history_response=$(curl -s -X GET "$API_BASE/v1/dialogs/$DIALOG_ID/messages?user_id=$TEST_USER")
        if echo "$history_response" | grep -q "messages"; then
            print_success "История диалога получена"
        else
            print_error "Ошибка получения истории диалога"
        fi
    else
        print_error "Ошибка создания диалога: $response"
        return 1
    fi
}

# Функция для тестирования создания пайплайна
test_pipeline_creation() {
    print_info "Тестирование создания пайплайна..."
    
    response=$(curl -s -X POST "$API_BASE/v1/pipelines" \
        -H "Content-Type: application/json" \
        -d "{
            \"source\": {
                \"type\": {
                    \"file\": {
                        \"format\": \"csv\",
                        \"url\": \"users/$TEST_USER/test-file.csv\"
                    }
                }
            },
            \"target\": {
                \"type\": \"postgres\",
                \"table_name\": \"test_table\"
            },
            \"user_id\": \"$TEST_USER\",
            \"file_id\": \"test-file-id\"
        }")
    
    if echo "$response" | grep -q "pipeline_id"; then
        PIPELINE_ID=$(echo "$response" | grep -o '"pipeline_id":"[^"]*"' | cut -d'"' -f4)
        print_success "Пайплайн создан, ID: $PIPELINE_ID"
        
        # Проверяем статус пайплайна
        status_response=$(curl -s -X GET "$API_BASE/v1/pipelines/$PIPELINE_ID")
        if echo "$status_response" | grep -q "pipeline_id"; then
            print_success "Статус пайплайна получен"
        else
            print_error "Ошибка получения статуса пайплайна"
        fi
    else
        print_error "Ошибка создания пайплайна: $response"
        return 1
    fi
}

# Функция для тестирования Airflow DAGs
test_airflow_dags() {
    print_info "Тестирование Airflow DAGs..."
    
    response=$(curl -s -X GET "$AIRFLOW_BASE/api/v1/dags")
    if echo "$response" | grep -q "dags"; then
        print_success "DAG в Airflow доступны"
    else
        print_error "Ошибка получения DAG из Airflow"
        return 1
    fi
}

# Функция для тестирования MinIO
test_minio() {
    print_info "Тестирование MinIO..."
    
    if docker-compose exec -T minio mc ls minio/files >/dev/null 2>&1; then
        print_success "MinIO доступен"
    else
        print_error "MinIO недоступен"
        return 1
    fi
}

# Функция для тестирования PostgreSQL
test_postgresql() {
    print_info "Тестирование PostgreSQL..."
    
    if docker-compose exec -T postgres pg_isready -U testuser >/dev/null 2>&1; then
        print_success "PostgreSQL доступна"
    else
        print_error "PostgreSQL недоступна"
        return 1
    fi
}

# Функция для тестирования Ollama
test_ollama() {
    print_info "Тестирование Ollama..."
    
    if curl -s http://localhost:11434/api/tags >/dev/null 2>&1; then
        print_success "Ollama доступен"
    else
        print_error "Ollama недоступен"
        return 1
    fi
}

# Основная функция
main() {
    print_header
    
    local failed_tests=0
    
    # Тестируем основные сервисы
    test_api_gateway || failed_tests=$((failed_tests + 1))
    test_file_service || failed_tests=$((failed_tests + 1))
    test_airflow || failed_tests=$((failed_tests + 1))
    
    echo ""
    
    # Тестируем инфраструктуру
    test_minio || failed_tests=$((failed_tests + 1))
    test_postgresql || failed_tests=$((failed_tests + 1))
    test_ollama || failed_tests=$((failed_tests + 1))
    
    echo ""
    
    # Тестируем функциональность
    test_small_file_upload || failed_tests=$((failed_tests + 1))
    test_dialog_creation || failed_tests=$((failed_tests + 1))
    test_pipeline_creation || failed_tests=$((failed_tests + 1))
    test_airflow_dags || failed_tests=$((failed_tests + 1))
    
    echo ""
    
    # Результаты
    if [ $failed_tests -eq 0 ]; then
        echo -e "${GREEN}🎉 Все тесты пройдены успешно!${NC}"
    else
        echo -e "${RED}❌ $failed_tests тестов не пройдены${NC}"
    fi
    
    echo ""
    echo "Для подробного тестирования используйте: ./test-workflow.sh"
}

# Запуск основной функции
main "$@"