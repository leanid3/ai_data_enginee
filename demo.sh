#!/bin/bash

# Демонстрационный скрипт для системы "Инженер данных"

set -e

# Цвета для вывода
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Переменные
API_BASE="http://localhost:8080"
TEST_USER="demo-user-$(date +%s)"

print_header() {
    echo -e "${BLUE}"
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║                    ДЕМОНСТРАЦИЯ СИСТЕМЫ                     ║"
    echo "║                    'ИНЖЕНЕР ДАННЫХ'                         ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

print_step() {
    echo -e "${PURPLE}🔄 Шаг $1: $2${NC}"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_info() {
    echo -e "${CYAN}ℹ${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Создание тестового файла
create_test_file() {
    print_info "Создание тестового файла..."
    
    cat > /tmp/demo-data.csv << EOF
id,name,email,age,department,salary
1,John Doe,john@example.com,30,Engineering,75000
2,Jane Smith,jane@example.com,28,Marketing,65000
3,Bob Johnson,bob@example.com,35,Sales,70000
4,Alice Brown,alice@example.com,32,Engineering,80000
5,Charlie Wilson,charlie@example.com,29,Marketing,60000
6,Diana Davis,diana@example.com,31,Sales,72000
7,Frank Miller,frank@example.com,33,Engineering,85000
8,Grace Lee,grace@example.com,27,Marketing,58000
9,Henry Taylor,henry@example.com,34,Sales,78000
10,Iris White,iris@example.com,26,Engineering,70000
EOF
    
    print_success "Тестовый файл создан: /tmp/demo-data.csv"
}

# Демонстрация загрузки файла
demo_file_upload() {
    print_step "1" "Загрузка файла в систему"
    
    response=$(curl -s -X POST "$API_BASE/v1/files/upload/csv" \
        -H "Content-Type: multipart/form-data" \
        -F "file=@/tmp/demo-data.csv" \
        -F "user_id=$TEST_USER")
    
    if echo "$response" | grep -q "file_id"; then
        FILE_ID=$(echo "$response" | grep -o '"file_id":"[^"]*"' | cut -d'"' -f4)
        print_success "Файл загружен успешно! ID: $FILE_ID"
        
        # Показываем информацию о файле
        print_info "Получение информации о файле..."
        file_info=$(curl -s -X GET "$API_BASE/v1/files/$FILE_ID?user_id=$TEST_USER")
        echo "$file_info" | jq '.' 2>/dev/null || echo "$file_info"
    else
        print_error "Ошибка загрузки файла: $response"
        return 1
    fi
}

# Демонстрация создания диалога
demo_dialog_creation() {
    print_step "2" "Создание диалога с системой"
    
    response=$(curl -s -X POST "$API_BASE/v1/dialogs" \
        -H "Content-Type: application/json" \
        -d "{
            \"user_id\": \"$TEST_USER\",
            \"title\": \"Анализ данных сотрудников\",
            \"initial_message\": \"Привет! Я загрузил файл с данными о сотрудниках. Можешь проанализировать его и дать рекомендации?\"
        }")
    
    if echo "$response" | grep -q "dialog_id"; then
        DIALOG_ID=$(echo "$response" | grep -o '"dialog_id":"[^"]*"' | cut -d'"' -f4)
        print_success "Диалог создан! ID: $DIALOG_ID"
        
        # Показываем историю диалога
        print_info "Получение истории диалога..."
        dialog_history=$(curl -s -X GET "$API_BASE/v1/dialogs/$DIALOG_ID/messages?user_id=$TEST_USER")
        echo "$dialog_history" | jq '.' 2>/dev/null || echo "$dialog_history"
    else
        print_error "Ошибка создания диалога: $response"
        return 1
    fi
}

# Демонстрация создания пайплайна
demo_pipeline_creation() {
    print_step "3" "Создание пайплайна анализа данных"
    
    response=$(curl -s -X POST "$API_BASE/v1/pipelines" \
        -H "Content-Type: application/json" \
        -d "{
            \"source\": {
                \"type\": {
                    \"file\": {
                        \"format\": \"csv\",
                        \"url\": \"users/$TEST_USER/$FILE_ID/demo-data.csv\"
                    }
                }
            },
            \"target\": {
                \"type\": \"postgres\",
                \"table_name\": \"employees\"
            },
            \"user_id\": \"$TEST_USER\",
            \"file_id\": \"$FILE_ID\"
        }")
    
    if echo "$response" | grep -q "pipeline_id"; then
        PIPELINE_ID=$(echo "$response" | grep -o '"pipeline_id":"[^"]*"' | cut -d'"' -f4)
        print_success "Пайплайн создан! ID: $PIPELINE_ID"
        
        # Показываем информацию о пайплайне
        print_info "Получение информации о пайплайне..."
        pipeline_info=$(curl -s -X GET "$API_BASE/v1/pipelines/$PIPELINE_ID")
        echo "$pipeline_info" | jq '.' 2>/dev/null || echo "$pipeline_info"
    else
        print_error "Ошибка создания пайплайна: $response"
        return 1
    fi
}

# Демонстрация проверки Airflow
demo_airflow_check() {
    print_step "4" "Проверка Airflow DAG"
    
    print_info "Получение списка DAG из Airflow..."
    dags_response=$(curl -s -X GET "http://localhost:8082/api/v1/dags")
    if echo "$dags_response" | grep -q "dags"; then
        print_success "DAG в Airflow доступны"
        echo "$dags_response" | jq '.dags[].dag_id' 2>/dev/null || echo "DAG найдены"
    else
        print_warning "Airflow недоступен или DAG не найдены"
    fi
}

# Демонстрация проверки MinIO
demo_minio_check() {
    print_step "5" "Проверка MinIO хранилища"
    
    print_info "Проверка файлов в MinIO..."
    if docker-compose exec -T minio mc ls minio/files >/dev/null 2>&1; then
        print_success "MinIO доступен"
        print_info "Список файлов в MinIO:"
        docker-compose exec -T minio mc ls minio/files/users/ 2>/dev/null || echo "Файлы пользователей не найдены"
    else
        print_warning "MinIO недоступен"
    fi
}

# Демонстрация проверки базы данных
demo_database_check() {
    print_step "6" "Проверка базы данных"
    
    print_info "Проверка подключения к PostgreSQL..."
    if docker-compose exec -T postgres pg_isready -U testuser >/dev/null 2>&1; then
        print_success "PostgreSQL доступна"
        
        print_info "Проверка таблиц в базе данных..."
        docker-compose exec -T postgres psql -U testuser -d testdb -c "SELECT COUNT(*) FROM files;" 2>/dev/null || echo "Таблица files не найдена"
        docker-compose exec -T postgres psql -U testuser -d testdb -c "SELECT COUNT(*) FROM dialogs;" 2>/dev/null || echo "Таблица dialogs не найдена"
    else
        print_warning "PostgreSQL недоступна"
    fi
}

# Демонстрация проверки LLM
demo_llm_check() {
    print_step "7" "Проверка LLM сервиса"
    
    print_info "Проверка Ollama..."
    if curl -s http://localhost:11434/api/tags >/dev/null 2>&1; then
        print_success "Ollama доступен"
        
        print_info "Доступные модели:"
        curl -s http://localhost:11434/api/tags | jq '.models[].name' 2>/dev/null || echo "Модели не найдены"
    else
        print_warning "Ollama недоступен"
    fi
}

# Демонстрация мониторинга
demo_monitoring() {
    print_step "8" "Мониторинг системы"
    
    print_info "Статус контейнеров:"
    docker-compose ps
    
    print_info "Использование ресурсов:"
    docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}" | head -10
    
    print_info "Логи последних событий:"
    docker-compose logs --tail=5 api-gateway
}

# Очистка после демонстрации
cleanup() {
    print_info "Очистка после демонстрации..."
    
    # Удаляем тестовый файл
    rm -f /tmp/demo-data.csv
    
    print_success "Очистка завершена"
}

# Основная функция
main() {
    print_header
    
    print_info "Демонстрация системы 'Инженер данных'"
    print_info "Пользователь: $TEST_USER"
    echo ""
    
    # Проверяем, что система запущена
    if ! curl -s "$API_BASE/health" >/dev/null 2>&1; then
        print_error "Система не запущена. Запустите: ./start.sh start"
        exit 1
    fi
    
    # Создаем тестовый файл
    create_test_file
    echo ""
    
    # Демонстрируем функциональность
    demo_file_upload
    echo ""
    
    demo_dialog_creation
    echo ""
    
    demo_pipeline_creation
    echo ""
    
    demo_airflow_check
    echo ""
    
    demo_minio_check
    echo ""
    
    demo_database_check
    echo ""
    
    demo_llm_check
    echo ""
    
    demo_monitoring
    echo ""
    
    # Очистка
    cleanup
    echo ""
    
    print_success "🎉 Демонстрация завершена!"
    echo ""
    echo "Для полного тестирования используйте:"
    echo "  ./quick-test.sh      - Быстрое тестирование"
    echo "  ./test-workflow.sh   - Полное тестирование"
    echo "  ./check-system.sh    - Проверка системы"
    echo ""
    echo "Документация:"
    echo "  README.md                    - Основная документация"
    echo "  TESTING_INSTRUCTIONS.md     - Инструкции по тестированию"
    echo "  USER_WORKFLOW_TESTING.md    - Тестирование workflow"
}

# Запуск основной функции
main "$@"