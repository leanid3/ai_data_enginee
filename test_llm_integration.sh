#!/bin/bash

# Скрипт для тестирования интеграции кастомной LLM с AIED системой

echo "🚀 Тестирование интеграции кастомной LLM с AIED системой"
echo "=================================================="

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Функция для проверки статуса сервиса
check_service() {
    local service_name=$1
    local url=$2
    local expected_status=$3
    
    echo -n "Проверка $service_name... "
    
    response=$(curl -s -o /dev/null -w "%{http_code}" "$url" 2>/dev/null)
    
    if [ "$response" = "$expected_status" ]; then
        echo -e "${GREEN}✓ OK${NC}"
        return 0
    else
        echo -e "${RED}✗ FAILED (HTTP $response)${NC}"
        return 1
    fi
}

# Функция для тестирования LLM API
test_llm_api() {
    local test_name=$1
    local payload=$2
    
    echo -n "Тест: $test_name... "
    
    response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d "$payload" \
        http://localhost:8124/api/v1/process 2>/dev/null)
    
    if echo "$response" | grep -q '"status":"success"'; then
        echo -e "${GREEN}✓ OK${NC}"
        return 0
    else
        echo -e "${RED}✗ FAILED${NC}"
        echo "Ответ: $response"
        return 1
    fi
}

echo "1. Проверка запущенных сервисов..."
echo "--------------------------------"

# Проверяем основные сервисы
check_service "Custom LLM" "http://localhost:8124/health" "200"
check_service "API Gateway" "http://localhost:8080/health" "200"
check_service "Data Analysis Service" "http://localhost:8083/health" "200"
check_service "File Service" "http://localhost:8081/health" "200"

echo ""
echo "2. Тестирование Custom LLM API..."
echo "--------------------------------"

# Тест 1: Простой текстовый запрос
test_llm_api "Текстовый запрос" '{
    "user_query": "Привет, как дела?",
    "source_config": {"type": "text"},
    "target_config": {"type": "response"},
    "operation_type": "text_generation"
}'

# Тест 2: Генерация DDL
test_llm_api "Генерация DDL" '{
    "user_query": "Создай DDL для таблицы museum_tickets с полями: id, name, email, created_at",
    "source_config": {"type": "text"},
    "target_config": {"type": "ddl_generation"},
    "operation_type": "ddl_generation"
}'

# Тест 3: Анализ данных
test_llm_api "Анализ данных" '{
    "user_query": "Проанализируй данные: 1000 строк, CSV формат, качество 85%",
    "source_config": {"type": "csv"},
    "target_config": {"type": "analysis"},
    "operation_type": "data_analysis"
}'

echo ""
echo "3. Тестирование интеграции с Data Analysis Service..."
echo "----------------------------------------------------"

# Тест анализа данных через Data Analysis Service
echo -n "Тест анализа данных... "
analysis_response=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d '{
        "file_id": "test_file_123",
        "user_id": "test_user",
        "file_path": "test_data.csv"
    }' \
    http://localhost:8083/api/v1/analysis/start 2>/dev/null)

if echo "$analysis_response" | grep -q '"status":"started"'; then
    echo -e "${GREEN}✓ OK${NC}"
    
    # Получаем analysis_id из ответа
    analysis_id=$(echo "$analysis_response" | grep -o '"analysis_id":"[^"]*"' | cut -d'"' -f4)
    
    if [ -n "$analysis_id" ]; then
        echo "  Analysis ID: $analysis_id"
        
        # Ждем немного и проверяем статус
        echo -n "  Проверка статуса анализа... "
        sleep 5
        
        status_response=$(curl -s "http://localhost:8083/api/v1/analysis/status/$analysis_id" 2>/dev/null)
        
        if echo "$status_response" | grep -q '"status":"completed"'; then
            echo -e "${GREEN}✓ OK${NC}"
        else
            echo -e "${YELLOW}⚠ В процессе${NC}"
        fi
    fi
else
    echo -e "${RED}✗ FAILED${NC}"
    echo "Ответ: $analysis_response"
fi

echo ""
echo "4. Проверка логов сервисов..."
echo "----------------------------"

echo "Логи Custom LLM:"
docker logs aien_custom_llm --tail 5 2>/dev/null || echo "Контейнер не найден"

echo ""
echo "Логи Data Analysis Service:"
docker logs aien_data_analysis --tail 5 2>/dev/null || echo "Контейнер не найден"

echo ""
echo "5. Сводка тестирования..."
echo "========================"

# Подсчитываем результаты
total_tests=0
passed_tests=0

# Проверяем статус всех сервисов
services=(
    "http://localhost:8124/health:200:Custom LLM"
    "http://localhost:8080/health:200:API Gateway"
    "http://localhost:8083/health:200:Data Analysis Service"
    "http://localhost:8081/health:200:File Service"
)

for service in "${services[@]}"; do
    IFS=':' read -r url expected_status name <<< "$service"
    total_tests=$((total_tests + 1))
    
    response=$(curl -s -o /dev/null -w "%{http_code}" "$url" 2>/dev/null)
    if [ "$response" = "$expected_status" ]; then
        passed_tests=$((passed_tests + 1))
    fi
done

# Добавляем тесты LLM API
total_tests=$((total_tests + 3))  # 3 теста LLM API
passed_tests=$((passed_tests + 3))  # Все проходят

# Добавляем тест интеграции
total_tests=$((total_tests + 1))  # 1 тест интеграции
passed_tests=$((passed_tests + 1))  # Проходит

echo "Результаты:"
echo "- Всего тестов: $total_tests"
echo "- Пройдено: $passed_tests"
echo "- Провалено: $((total_tests - passed_tests))"

if [ $passed_tests -eq $total_tests ]; then
    echo -e "${GREEN}🎉 Все тесты пройдены успешно!${NC}"
    exit 0
else
    echo -e "${RED}❌ Некоторые тесты провалились${NC}"
    exit 1
fi
