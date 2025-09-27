#!/bin/bash

# AIEN Backend - Запуск приложения
# Автор: AIEN Team
# Версия: 1.0.0

set -e

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Функция для вывода сообщений
print_message() {
    echo -e "${BLUE}[AIEN]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Проверка зависимостей
check_dependencies() {
    print_message "Проверка зависимостей..."
    
    if ! command -v docker &> /dev/null; then
        print_error "Docker не установлен. Установите Docker: https://docs.docker.com/get-docker/"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose не установлен. Установите Docker Compose: https://docs.docker.com/compose/install/"
        exit 1
    fi
    
    print_success "Все зависимости установлены"
}

# Проверка портов
check_ports() {
    print_message "Проверка доступности портов..."
    
    local ports=(8080 8081 8083 8084 5432 6379 9000 9001 11434)
    local occupied_ports=()
    
    for port in "${ports[@]}"; do
        if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
            occupied_ports+=($port)
        fi
    done
    
    if [ ${#occupied_ports[@]} -gt 0 ]; then
        print_warning "Следующие порты заняты: ${occupied_ports[*]}"
        print_warning "Возможно, приложение уже запущено или порты используются другими сервисами"
        read -p "Продолжить? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
    
    print_success "Порты свободны"
}

# Создание необходимых директорий
create_directories() {
    print_message "Создание директорий..."
    
    mkdir -p logs
    mkdir -p data/postgres
    mkdir -p data/minio
    mkdir -p data/ollama
    
    print_success "Директории созданы"
}

# Сборка образов
build_images() {
    print_message "Сборка Docker образов..."
    
    # Сборка всех сервисов
    docker-compose -f docker-compose.prod.yml build --parallel
    
    print_success "Образы собраны"
}

# Запуск сервисов
start_services() {
    print_message "Запуск сервисов..."
    
    # Запуск в фоновом режиме
    docker-compose -f docker-compose.prod.yml up -d
    
    print_success "Сервисы запущены"
}

# Ожидание готовности сервисов
wait_for_services() {
    print_message "Ожидание готовности сервисов..."
    
    local max_attempts=60
    local attempt=0
    
    while [ $attempt -lt $max_attempts ]; do
        if curl -s http://localhost:8080/health >/dev/null 2>&1; then
            print_success "API Gateway готов"
            break
        fi
        
        attempt=$((attempt + 1))
        print_message "Попытка $attempt/$max_attempts - ожидание API Gateway..."
        sleep 5
    done
    
    if [ $attempt -eq $max_attempts ]; then
        print_error "API Gateway не отвечает после $max_attempts попыток"
        return 1
    fi
    
    # Проверка других сервисов
    sleep 10
    
    if curl -s http://localhost:8083/health >/dev/null 2>&1; then
        print_success "Data Analysis Service готов"
    else
        print_warning "Data Analysis Service не отвечает"
    fi
    
    if curl -s http://localhost:8081/health >/dev/null 2>&1; then
        print_success "File Service готов"
    else
        print_warning "File Service не отвечает"
    fi
}

# Показать статус сервисов
show_status() {
    print_message "Статус сервисов:"
    echo
    
    # API Gateway
    if curl -s http://localhost:8080/health >/dev/null 2>&1; then
        print_success "✅ API Gateway: http://localhost:8080"
    else
        print_error "❌ API Gateway: недоступен"
    fi
    
    # File Service
    if curl -s http://localhost:8081/health >/dev/null 2>&1; then
        print_success "✅ File Service: http://localhost:8081"
    else
        print_error "❌ File Service: недоступен"
    fi
    
    # Data Analysis Service
    if curl -s http://localhost:8083/health >/dev/null 2>&1; then
        print_success "✅ Data Analysis Service: http://localhost:8083"
    else
        print_error "❌ Data Analysis Service: недоступен"
    fi
    
    # Adminer
    if curl -s http://localhost:8084 >/dev/null 2>&1; then
        print_success "✅ Adminer: http://localhost:8084"
    else
        print_error "❌ Adminer: недоступен"
    fi
    
    # MinIO
    if curl -s http://localhost:9000 >/dev/null 2>&1; then
        print_success "✅ MinIO: http://localhost:9000"
    else
        print_error "❌ MinIO: недоступен"
    fi
    
    echo
    print_message "Полезные ссылки:"
    echo "  📊 API Gateway: http://localhost:8080"
    echo "  📁 File Service: http://localhost:8081"
    echo "  🔍 Data Analysis: http://localhost:8083"
    echo "  🗄️  Database Admin: http://localhost:8084"
    echo "  📦 MinIO Console: http://localhost:9001"
    echo "  🤖 Ollama: http://localhost:11434"
}

# Показать логи
show_logs() {
    print_message "Логи сервисов (Ctrl+C для выхода):"
    docker-compose -f docker-compose.prod.yml logs -f
}

# Остановка сервисов
stop_services() {
    print_message "Остановка сервисов..."
    docker-compose -f docker-compose.prod.yml down
    print_success "Сервисы остановлены"
}

# Очистка данных
clean_data() {
    print_warning "Это удалит ВСЕ данные (база данных, файлы, модели)"
    read -p "Продолжить? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_message "Остановка сервисов..."
        docker-compose -f docker-compose.prod.yml down -v
        print_message "Удаление данных..."
        sudo rm -rf data/
        print_success "Данные очищены"
    fi
}

# Показать помощь
show_help() {
    echo "AIEN Backend - Система анализа данных"
    echo
    echo "Использование: $0 [команда]"
    echo
    echo "Команды:"
    echo "  start     - Запустить приложение (по умолчанию)"
    echo "  stop      - Остановить приложение"
    echo "  restart   - Перезапустить приложение"
    echo "  status    - Показать статус сервисов"
    echo "  logs      - Показать логи"
    echo "  clean     - Очистить все данные"
    echo "  help      - Показать эту справку"
    echo
    echo "Примеры:"
    echo "  $0 start    # Запустить приложение"
    echo "  $0 status   # Проверить статус"
    echo "  $0 logs     # Посмотреть логи"
}

# Основная функция
main() {
    echo "🚀 AIEN Backend - Система анализа данных"
    echo "========================================"
    echo
    
    case "${1:-start}" in
        "start")
            check_dependencies
            check_ports
            create_directories
            build_images
            start_services
            wait_for_services
            show_status
            ;;
        "stop")
            stop_services
            ;;
        "restart")
            stop_services
            sleep 5
            main start
            ;;
        "status")
            show_status
            ;;
        "logs")
            show_logs
            ;;
        "clean")
            clean_data
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            print_error "Неизвестная команда: $1"
            show_help
            exit 1
            ;;
    esac
}

# Запуск
main "$@"