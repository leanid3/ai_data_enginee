# 🚀 AIEN Backend - Система анализа данных

## 📋 Описание

AIEN Backend - это микросервисная система для анализа данных, которая использует LLM (Ollama) для генерации DDL скриптов и рекомендаций по хранению данных. Система поддерживает загрузку файлов в различных форматах (CSV, JSON, XML) и автоматически анализирует их структуру.

## 🏗️ Архитектура

### Микросервисы:
- **API Gateway** (8080) - точка входа для всех запросов
- **File Service** (8081) - загрузка и управление файлами
- **Data Analysis Service** (8083) - анализ данных с помощью LLM
- **Data Profiler** (50052) - профилирование данных
- **Data Transfer** (50053) - перенос данных в целевые системы
- **Chat Service** (50055) - чат-интерфейс
- **LLM Service** (50056) - работа с языковыми моделями
- **Orchestration Service** (50057) - оркестрация процессов

### Внешние сервисы:
- **PostgreSQL** (5432) - основная база данных
- **Redis** (6379) - кэширование
- **MinIO** (9000/9001) - объектное хранилище
- **Ollama** (11434) - LLM для анализа данных
- **Adminer** (8084) - администрирование БД

## 🚀 Быстрый запуск

### Требования:
- Docker 20.10+
- Docker Compose 2.0+
- 8GB RAM (рекомендуется)
- 20GB свободного места

### Установка:

1. **Клонируйте репозиторий:**
```bash
git clone https://github.com/your-username/AIEN_backend.git
cd AIEN_backend
```

2. **Запустите приложение:**
```bash
./start.sh
```

3. **Проверьте статус:**
```bash
./start.sh status
```

## 📊 Использование

### API Endpoints:

#### Health Checks:
```bash
curl http://localhost:8080/health
curl http://localhost:8081/health
curl http://localhost:8083/health
```

#### Загрузка файла:
```bash
curl -X POST -F "file=@data.csv" "http://localhost:8081/api/v1/files/upload?user_id=test-user"
```

#### Анализ данных:
```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"file_id":"file-id","user_id":"user-id","file_path":"path/to/file"}' \
  http://localhost:8083/api/v1/analysis/start
```

#### Статус анализа:
```bash
curl http://localhost:8083/api/v1/analysis/status/analysis-id
```

## 🧪 Тестирование

### Postman коллекция:
1. Импортируйте `AIEN_Backend_API.postman_collection.json` в Postman
2. Настройте переменные для ваших файлов
3. Запустите тестовые сценарии

### Тестовые файлы:
- `test_data/sales_data.csv` - данные продаж
- `test_data/users.json` - пользовательские данные
- `test_data/products.xml` - каталог товаров

## 🛠️ Управление

### Команды скрипта:

```bash
./start.sh start     # Запустить приложение
./start.sh stop      # Остановить приложение
./start.sh restart   # Перезапустить приложение
./start.sh status    # Показать статус сервисов
./start.sh logs      # Показать логи
./start.sh clean     # Очистить все данные
./start.sh help      # Показать справку
```

### Полезные ссылки:
- **API Gateway:** http://localhost:8080
- **File Service:** http://localhost:8081
- **Data Analysis:** http://localhost:8083
- **Database Admin:** http://localhost:8084
- **MinIO Console:** http://localhost:9001
- **Ollama:** http://localhost:11434

## 🔧 Конфигурация

### Переменные окружения:
- `POSTGRES_*` - настройки PostgreSQL
- `REDIS_*` - настройки Redis
- `MINIO_*` - настройки MinIO
- `OLLAMA_URL` - URL Ollama сервиса

### Порты:
- 8080 - API Gateway
- 8081 - File Service
- 8083 - Data Analysis Service
- 8084 - Adminer
- 5432 - PostgreSQL
- 6379 - Redis
- 9000/9001 - MinIO
- 11434 - Ollama

## 📈 Мониторинг

### Логи:
```bash
./start.sh logs
```

### Статус сервисов:
```bash
./start.sh status
```

### Docker контейнеры:
```bash
docker-compose -f docker-compose.prod.yml ps
```

## 🗄️ Данные

### Персистентные данные:
- `data/postgres/` - база данных PostgreSQL
- `data/minio/` - файлы MinIO
- `data/ollama/` - модели Ollama

### Очистка данных:
```bash
./start.sh clean
```

## 🚨 Устранение неполадок

### Проблемы с портами:
```bash
# Проверить занятые порты
lsof -i :8080
lsof -i :8081
lsof -i :8083

# Освободить порты
sudo kill -9 $(lsof -t -i:8080)
```

### Проблемы с Docker:
```bash
# Перезапустить Docker
sudo systemctl restart docker

# Очистить Docker
docker system prune -a
```

### Проблемы с памятью:
```bash
# Проверить использование памяти
docker stats

# Ограничить память для Ollama
docker-compose -f docker-compose.prod.yml up -d --scale ollama=0
```

## 📝 Разработка

### Структура проекта:
```
AIEN_backend/
├── api-geteway/          # API Gateway
├── data-profiler/         # Data Profiler
├── data-transfer/        # Data Transfer
├── file-service/         # File Service
├── data-analysis-service/ # Data Analysis
├── chat-service/         # Chat Service
├── llm-service/          # LLM Service
├── orchestration-service/ # Orchestration
├── test_data/           # Тестовые файлы
├── docker-compose.prod.yml # Production конфигурация
├── start.sh             # Скрипт запуска
└── README.md            # Документация
```

### Добавление нового сервиса:
1. Создайте директорию сервиса
2. Добавьте Dockerfile
3. Обновите docker-compose.prod.yml
4. Добавьте в start.sh

## 🤝 Вклад в проект

1. Fork репозитория
2. Создайте feature branch
3. Внесите изменения
4. Создайте Pull Request

## 📄 Лицензия

MIT License - см. файл LICENSE

## 🆘 Поддержка

- **Issues:** https://github.com/your-username/AIEN_backend/issues
- **Discussions:** https://github.com/your-username/AIEN_backend/discussions
- **Email:** support@aien.ai

---

**AIEN Backend** - Система анализа данных с поддержкой LLM 🚀