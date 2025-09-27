#!/bin/bash

# Простое тестирование workflow

set -e

# Цвета для вывода
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Переменные
TEST_USER="test-user-$(date +%s)"
CSV_FILE="_data/csv/part-00000-37dced01-2ad2-48c8-a56d-54b4d8760599-c000.csv"

print_header() {
    echo -e "${BLUE}"
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║              ТЕСТИРОВАНИЕ ПРОСТОГО WORKFLOW                  ║"
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

# Тест 1: Загрузка файла
test_file_upload() {
    print_info "Тест 1: Загрузка CSV файла..."
    
    if [ ! -f "$CSV_FILE" ]; then
        print_error "CSV файл не найден: $CSV_FILE"
        return 1
    fi
    
    response=$(curl -s -X POST http://localhost:50054/v1/files/upload/csv \
        -H "Content-Type: multipart/form-data" \
        -F "file=@$CSV_FILE" \
        -F "user_id=$TEST_USER")
    
    if echo "$response" | grep -q "file_id"; then
        FILE_ID=$(echo "$response" | grep -o '"file_id":"[^"]*"' | cut -d'"' -f4)
        print_success "CSV файл загружен, ID: $FILE_ID"
        echo "Ответ: $response"
    else
        print_error "Ошибка загрузки CSV файла: $response"
        return 1
    fi
}

# Тест 2: Создание диалога
test_create_dialog() {
    print_info "Тест 2: Создание диалога..."
    
    response=$(curl -s -X POST http://localhost:50055/v1/dialogs \
        -H "Content-Type: application/json" \
        -d "{
            \"user_id\": \"$TEST_USER\",
            \"title\": \"Анализ данных\",
            \"initial_message\": \"Проанализируй загруженный файл\"
        }")
    
    if echo "$response" | grep -q "dialog_id"; then
        DIALOG_ID=$(echo "$response" | grep -o '"dialog_id":"[^"]*"' | cut -d'"' -f4)
        print_success "Диалог создан, ID: $DIALOG_ID"
        echo "Ответ: $response"
    else
        print_error "Ошибка создания диалога: $response"
        return 1
    fi
}

# Тест 3: Анализ данных
test_data_analysis() {
    print_info "Тест 3: Анализ данных через LLM..."
    
    response=$(curl -s -X POST http://localhost:50056/v1/analyze \
        -H "Content-Type: application/json" \
        -d "{
            \"user_id\": \"$TEST_USER\",
            \"file_path\": \"/data/test.csv\",
            \"file_format\": \"csv\",
            \"sample_size\": 1000,
            \"analysis_type\": \"detailed\"
        }")
    
    if echo "$response" | grep -q "request_id"; then
        REQUEST_ID=$(echo "$response" | grep -o '"request_id":"[^"]*"' | cut -d'"' -f4)
        print_success "Анализ данных запущен, ID: $REQUEST_ID"
        echo "Ответ: $response"
    else
        print_error "Ошибка анализа данных: $response"
        return 1
    fi
}

# Тест 4: Создание пайплайна
test_create_pipeline() {
    print_info "Тест 4: Создание пайплайна..."
    
    response=$(curl -s -X POST http://localhost:50057/v1/dags \
        -H "Content-Type: application/json" \
        -d "{
            \"user_id\": \"$TEST_USER\",
            \"dag_id\": \"data_analysis_$TEST_USER\",
            \"dag_yaml\": \"test yaml content\",
            \"schedule_interval\": \"manual\",
            \"start_immediately\": true,
            \"description\": \"Пайплайн анализа данных\"
        }")
    
    if echo "$response" | grep -q "dag_id"; then
        DAG_ID=$(echo "$response" | grep -o '"dag_id":"[^"]*"' | cut -d'"' -f4)
        print_success "Пайплайн создан, ID: $DAG_ID"
        echo "Ответ: $response"
    else
        print_error "Ошибка создания пайплайна: $response"
        return 1
    fi
}

# Тест 5: Проверка статуса пайплайна
test_pipeline_status() {
    print_info "Тест 5: Проверка статуса пайплайна..."
    
    response=$(curl -s -X GET "http://localhost:50057/v1/dags/data_analysis_$TEST_USER")
    
    if echo "$response" | grep -q "dag_id"; then
        print_success "Статус пайплайна получен"
        echo "Ответ: $response"
    else
        print_error "Ошибка получения статуса пайплайна: $response"
        return 1
    fi
}

# Тест 6: Чат с LLM
test_llm_chat() {
    print_info "Тест 6: Чат с LLM..."
    
    response=$(curl -s -X POST http://localhost:50056/v1/chat \
        -H "Content-Type: application/json" \
        -d "{
            \"user_id\": \"$TEST_USER\",
            \"message\": \"Проанализируй данные и дай рекомендации\",
            \"context\": \"Анализ CSV файла\"
        }")
    
    if echo "$response" | grep -q "request_id"; then
        print_success "Чат с LLM успешен"
        echo "Ответ: $response"
    else
        print_error "Ошибка чата с LLM: $response"
        return 1
    fi
}

# Основная функция
main() {
    print_header
    
    local failed_tests=0
    
    test_file_upload || failed_tests=$((failed_tests + 1))
    echo ""
    
    test_create_dialog || failed_tests=$((failed_tests + 1))
    echo ""
    
    test_data_analysis || failed_tests=$((failed_tests + 1))
    echo ""
    
    test_create_pipeline || failed_tests=$((failed_tests + 1))
    echo ""
    
    test_pipeline_status || failed_tests=$((failed_tests + 1))
    echo ""
    
    test_llm_chat || failed_tests=$((failed_tests + 1))
    echo ""
    
    # Результаты
    if [ $failed_tests -eq 0 ]; then
        echo -e "${GREEN}🎉 Все тесты пройдены успешно!${NC}"
        echo ""
        echo "Результаты тестирования:"
        echo "  - Пользователь: $TEST_USER"
        echo "  - Файл ID: $FILE_ID"
        echo "  - Диалог ID: $DIALOG_ID"
        echo "  - Запрос ID: $REQUEST_ID"
        echo "  - DAG ID: $DAG_ID"
    else
        echo -e "${RED}❌ $failed_tests тестов не пройдены${NC}"
    fi
}

# Запуск основной функции
main "$@"