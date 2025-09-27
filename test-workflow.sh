#!/bin/bash

# Скрипт для тестирования пользовательского workflow

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
TEST_USER="test-user-$(date +%s)"
CSV_FILE="_data/csv/part-00000-37dced01-2ad2-48c8-a56d-54b4d8760599-c000.csv"
JSON_FILE="_data/json/part1.json"
XML_FILE="_data/xml/part1.xml"

# Результаты тестов
declare -A test_results

print_header() {
    echo -e "${BLUE}"
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║              ТЕСТИРОВАНИЕ ПОЛЬЗОВАТЕЛЬСКОГО WORKFLOW         ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
    test_results["$2"]="PASS"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
    test_results["$2"]="FAIL"
}

print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Проверка доступности сервисов
check_services() {
    print_info "Проверка доступности сервисов..."
    
    # API Gateway
    if curl -s "$API_BASE/health" >/dev/null 2>&1; then
        print_success "API Gateway доступен" "api_gateway"
    else
        print_error "API Gateway недоступен" "api_gateway"
        exit 1
    fi
    
    # File Service
    if curl -s "$FILE_SERVICE_BASE/health" >/dev/null 2>&1; then
        print_success "File Service доступен" "file_service"
    else
        print_error "File Service недоступен" "file_service"
        exit 1
    fi
    
    # Airflow
    if curl -s "$AIRFLOW_BASE/health" >/dev/null 2>&1; then
        print_success "Airflow доступен" "airflow"
    else
        print_error "Airflow недоступен" "airflow"
        exit 1
    fi
}

# Тест 1: Загрузка CSV файла
test_csv_upload() {
    print_info "Тест 1: Загрузка CSV файла..."
    
    if [ ! -f "$CSV_FILE" ]; then
        print_error "CSV файл не найден: $CSV_FILE" "csv_upload"
        return
    fi
    
    response=$(curl -s -X POST "$API_BASE/v1/files/upload/csv" \
        -H "Content-Type: multipart/form-data" \
        -F "file=@$CSV_FILE" \
        -F "user_id=$TEST_USER")
    
    if echo "$response" | grep -q "file_id"; then
        CSV_FILE_ID=$(echo "$response" | grep -o '"file_id":"[^"]*"' | cut -d'"' -f4)
        print_success "CSV файл загружен, ID: $CSV_FILE_ID" "csv_upload"
    else
        print_error "Ошибка загрузки CSV файла: $response" "csv_upload"
    fi
}

# Тест 2: Загрузка JSON файла
test_json_upload() {
    print_info "Тест 2: Загрузка JSON файла..."
    
    if [ ! -f "$JSON_FILE" ]; then
        print_error "JSON файл не найден: $JSON_FILE" "json_upload"
        return
    fi
    
    response=$(curl -s -X POST "$API_BASE/v1/files/upload/json" \
        -H "Content-Type: multipart/form-data" \
        -F "file=@$JSON_FILE" \
        -F "user_id=$TEST_USER")
    
    if echo "$response" | grep -q "file_id"; then
        JSON_FILE_ID=$(echo "$response" | grep -o '"file_id":"[^"]*"' | cut -d'"' -f4)
        print_success "JSON файл загружен, ID: $JSON_FILE_ID" "json_upload"
    else
        print_error "Ошибка загрузки JSON файла: $response" "json_upload"
    fi
}

# Тест 3: Загрузка XML файла
test_xml_upload() {
    print_info "Тест 3: Загрузка XML файла..."
    
    if [ ! -f "$XML_FILE" ]; then
        print_error "XML файл не найден: $XML_FILE" "xml_upload"
        return
    fi
    
    response=$(curl -s -X POST "$API_BASE/v1/files/upload/xml" \
        -H "Content-Type: multipart/form-data" \
        -F "file=@$XML_FILE" \
        -F "user_id=$TEST_USER")
    
    if echo "$response" | grep -q "file_id"; then
        XML_FILE_ID=$(echo "$response" | grep -o '"file_id":"[^"]*"' | cut -d'"' -f4)
        print_success "XML файл загружен, ID: $XML_FILE_ID" "xml_upload"
    else
        print_error "Ошибка загрузки XML файла: $response" "xml_upload"
    fi
}

# Тест 4: Создание диалога
test_create_dialog() {
    print_info "Тест 4: Создание диалога..."
    
    response=$(curl -s -X POST "$API_BASE/v1/dialogs" \
        -H "Content-Type: application/json" \
        -d "{
            \"user_id\": \"$TEST_USER\",
            \"title\": \"Анализ данных\",
            \"initial_message\": \"Проанализируй загруженные файлы\"
        }")
    
    if echo "$response" | grep -q "dialog_id"; then
        DIALOG_ID=$(echo "$response" | grep -o '"dialog_id":"[^"]*"' | cut -d'"' -f4)
        print_success "Диалог создан, ID: $DIALOG_ID" "create_dialog"
    else
        print_error "Ошибка создания диалога: $response" "create_dialog"
    fi
}

# Тест 5: Создание пайплайна для CSV
test_create_csv_pipeline() {
    print_info "Тест 5: Создание пайплайна для CSV..."
    
    if [ -z "$CSV_FILE_ID" ]; then
        print_error "CSV файл не загружен" "csv_pipeline"
        return
    fi
    
    response=$(curl -s -X POST "$API_BASE/v1/pipelines" \
        -H "Content-Type: application/json" \
        -d "{
            \"source\": {
                \"type\": {
                    \"file\": {
                        \"format\": \"csv\",
                        \"url\": \"users/$TEST_USER/$CSV_FILE_ID/part-00000-37dced01-2ad2-48c8-a56d-54b4d8760599-c000.csv\"
                    }
                }
            },
            \"target\": {
                \"type\": \"postgres\",
                \"table_name\": \"analyzed_data\"
            },
            \"user_id\": \"$TEST_USER\",
            \"file_id\": \"$CSV_FILE_ID\"
        }")
    
    if echo "$response" | grep -q "pipeline_id"; then
        CSV_PIPELINE_ID=$(echo "$response" | grep -o '"pipeline_id":"[^"]*"' | cut -d'"' -f4)
        print_success "Пайплайн для CSV создан, ID: $CSV_PIPELINE_ID" "csv_pipeline"
    else
        print_error "Ошибка создания пайплайна для CSV: $response" "csv_pipeline"
    fi
}

# Тест 6: Создание пайплайна для JSON
test_create_json_pipeline() {
    print_info "Тест 6: Создание пайплайна для JSON..."
    
    if [ -z "$JSON_FILE_ID" ]; then
        print_error "JSON файл не загружен" "json_pipeline"
        return
    fi
    
    response=$(curl -s -X POST "$API_BASE/v1/pipelines" \
        -H "Content-Type: application/json" \
        -d "{
            \"source\": {
                \"type\": {
                    \"file\": {
                        \"format\": \"json\",
                        \"url\": \"users/$TEST_USER/$JSON_FILE_ID/part1.json\"
                    }
                }
            },
            \"target\": {
                \"type\": \"clickhouse\",
                \"table_name\": \"json_data\"
            },
            \"user_id\": \"$TEST_USER\",
            \"file_id\": \"$JSON_FILE_ID\"
        }")
    
    if echo "$response" | grep -q "pipeline_id"; then
        JSON_PIPELINE_ID=$(echo "$response" | grep -o '"pipeline_id":"[^"]*"' | cut -d'"' -f4)
        print_success "Пайплайн для JSON создан, ID: $JSON_PIPELINE_ID" "json_pipeline"
    else
        print_error "Ошибка создания пайплайна для JSON: $response" "json_pipeline"
    fi
}

# Тест 7: Создание пайплайна для XML
test_create_xml_pipeline() {
    print_info "Тест 7: Создание пайплайна для XML..."
    
    if [ -z "$XML_FILE_ID" ]; then
        print_error "XML файл не загружен" "xml_pipeline"
        return
    fi
    
    response=$(curl -s -X POST "$API_BASE/v1/pipelines" \
        -H "Content-Type: application/json" \
        -d "{
            \"source\": {
                \"type\": {
                    \"file\": {
                        \"format\": \"xml\",
                        \"url\": \"users/$TEST_USER/$XML_FILE_ID/part1.xml\"
                    }
                }
            },
            \"target\": {
                \"type\": \"hdfs\",
                \"table_name\": \"/data/xml_processed\"
            },
            \"user_id\": \"$TEST_USER\",
            \"file_id\": \"$XML_FILE_ID\"
        }")
    
    if echo "$response" | grep -q "pipeline_id"; then
        XML_PIPELINE_ID=$(echo "$response" | grep -o '"pipeline_id":"[^"]*"' | cut -d'"' -f4)
        print_success "Пайплайн для XML создан, ID: $XML_PIPELINE_ID" "xml_pipeline"
    else
        print_error "Ошибка создания пайплайна для XML: $response" "xml_pipeline"
    fi
}

# Тест 8: Проверка статуса пайплайнов
test_pipeline_status() {
    print_info "Тест 8: Проверка статуса пайплайнов..."
    
    if [ -n "$CSV_PIPELINE_ID" ]; then
        response=$(curl -s -X GET "$API_BASE/v1/pipelines/$CSV_PIPELINE_ID")
        if echo "$response" | grep -q "pipeline_id"; then
            print_success "Статус пайплайна CSV получен" "csv_pipeline_status"
        else
            print_error "Ошибка получения статуса пайплайна CSV" "csv_pipeline_status"
        fi
    fi
    
    if [ -n "$JSON_PIPELINE_ID" ]; then
        response=$(curl -s -X GET "$API_BASE/v1/pipelines/$JSON_PIPELINE_ID")
        if echo "$response" | grep -q "pipeline_id"; then
            print_success "Статус пайплайна JSON получен" "json_pipeline_status"
        else
            print_error "Ошибка получения статуса пайплайна JSON" "json_pipeline_status"
        fi
    fi
    
    if [ -n "$XML_PIPELINE_ID" ]; then
        response=$(curl -s -X GET "$API_BASE/v1/pipelines/$XML_PIPELINE_ID")
        if echo "$response" | grep -q "pipeline_id"; then
            print_success "Статус пайплайна XML получен" "xml_pipeline_status"
        else
            print_error "Ошибка получения статуса пайплайна XML" "xml_pipeline_status"
        fi
    fi
}

# Тест 9: Проверка DAG в Airflow
test_airflow_dags() {
    print_info "Тест 9: Проверка DAG в Airflow..."
    
    response=$(curl -s -X GET "$AIRFLOW_BASE/api/v1/dags")
    if echo "$response" | grep -q "dags"; then
        print_success "DAG в Airflow доступны" "airflow_dags"
    else
        print_error "Ошибка получения DAG из Airflow" "airflow_dags"
    fi
}

# Тест 10: Проверка информации о файлах
test_file_info() {
    print_info "Тест 10: Проверка информации о файлах..."
    
    if [ -n "$CSV_FILE_ID" ]; then
        response=$(curl -s -X GET "$API_BASE/v1/files/$CSV_FILE_ID?user_id=$TEST_USER")
        if echo "$response" | grep -q "file_id"; then
            print_success "Информация о CSV файле получена" "csv_file_info"
        else
            print_error "Ошибка получения информации о CSV файле" "csv_file_info"
        fi
    fi
    
    if [ -n "$JSON_FILE_ID" ]; then
        response=$(curl -s -X GET "$API_BASE/v1/files/$JSON_FILE_ID?user_id=$TEST_USER")
        if echo "$response" | grep -q "file_id"; then
            print_success "Информация о JSON файле получена" "json_file_info"
        else
            print_error "Ошибка получения информации о JSON файле" "json_file_info"
        fi
    fi
    
    if [ -n "$XML_FILE_ID" ]; then
        response=$(curl -s -X GET "$API_BASE/v1/files/$XML_FILE_ID?user_id=$TEST_USER")
        if echo "$response" | grep -q "file_id"; then
            print_success "Информация о XML файле получена" "xml_file_info"
        else
            print_error "Ошибка получения информации о XML файле" "xml_file_info"
        fi
    fi
}

# Тест 11: Проверка истории диалога
test_dialog_history() {
    print_info "Тест 11: Проверка истории диалога..."
    
    if [ -n "$DIALOG_ID" ]; then
        response=$(curl -s -X GET "$API_BASE/v1/dialogs/$DIALOG_ID/messages?user_id=$TEST_USER")
        if echo "$response" | grep -q "messages"; then
            print_success "История диалога получена" "dialog_history"
        else
            print_error "Ошибка получения истории диалога" "dialog_history"
        fi
    fi
}

# Генерация отчета
generate_report() {
    print_info "Генерация отчета о тестировании..."
    
    echo ""
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║                        ОТЧЕТ О ТЕСТИРОВАНИИ                  ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo ""
    
    total_tests=0
    passed_tests=0
    
    for test_name in "${!test_results[@]}"; do
        total_tests=$((total_tests + 1))
        if [ "${test_results[$test_name]}" = "PASS" ]; then
            passed_tests=$((passed_tests + 1))
            echo -e "${GREEN}✓${NC} $test_name: PASS"
        else
            echo -e "${RED}✗${NC} $test_name: FAIL"
        fi
    done
    
    echo ""
    echo "Результаты: $passed_tests/$total_tests тестов пройдено"
    
    if [ $passed_tests -eq $total_tests ]; then
        echo -e "${GREEN}🎉 Все тесты пройдены успешно!${NC}"
    else
        echo -e "${RED}❌ Некоторые тесты не пройдены${NC}"
    fi
    
    echo ""
    echo "Тестовые данные:"
    echo "  - Пользователь: $TEST_USER"
    echo "  - CSV файл ID: $CSV_FILE_ID"
    echo "  - JSON файл ID: $JSON_FILE_ID"
    echo "  - XML файл ID: $XML_FILE_ID"
    echo "  - Диалог ID: $DIALOG_ID"
    echo "  - CSV пайплайн ID: $CSV_PIPELINE_ID"
    echo "  - JSON пайплайн ID: $JSON_PIPELINE_ID"
    echo "  - XML пайплайн ID: $XML_PIPELINE_ID"
}

# Основная функция
main() {
    print_header
    
    check_services
    echo ""
    
    test_csv_upload
    echo ""
    
    test_json_upload
    echo ""
    
    test_xml_upload
    echo ""
    
    test_create_dialog
    echo ""
    
    test_create_csv_pipeline
    echo ""
    
    test_create_json_pipeline
    echo ""
    
    test_create_xml_pipeline
    echo ""
    
    test_pipeline_status
    echo ""
    
    test_airflow_dags
    echo ""
    
    test_file_info
    echo ""
    
    test_dialog_history
    echo ""
    
    generate_report
}

# Запуск основной функции
main "$@"