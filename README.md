# AIED Backend - Система анализа данных с кастомной LLM

## Описание

AIED Backend - это система для анализа данных с использованием кастомной LLM модели. Система позволяет загружать файлы, анализировать их структуру и качество, генерировать рекомендации по выбору хранилища и создавать DDL скрипты.

## Архитектура системы

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │────│  AI Backend     │────│  Custom LLM    │
│   (Port 3001)   │    │   (Port 8080)   │    │  (Port 8124)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │                       │
                                │                       │
                       ┌─────────────────┐    ┌─────────────────┐
                       │   PostgreSQL    │    │     Redis       │
                       │   (Port 5432)   │    │   (Port 6379)   │
                       └─────────────────┘    └─────────────────┘
                                │
                       ┌─────────────────┐    ┌─────────────────┐
                       │     MinIO       │    │   ClickHouse    │
                       │   (Port 9000)   │    │   (Port 9000)   │
                       └─────────────────┘    └─────────────────┘
```

## Компоненты системы

### Frontend (React)
- **Порт:** 3001
- **Функции:** Загрузка файлов, отображение результатов анализа
- **API:** Взаимодействует с AI Backend

### AI Data Engineer Backend (Go)
- **Порт:** 8080
- **Функции:** Основной API сервис, обработка данных, интеграция с LLM
- **Endpoints:**
  - `POST /api/v1/files/upload` - Загрузка файлов
  - `POST /api/v1/analysis/start` - Запуск анализа
  - `GET /api/v1/analysis/status/{id}` - Статус анализа
  - `GET /api/v1/health` - Проверка здоровья сервиса

### Custom LLM Service (Python)
- **Порт:** 8124
- **Функции:** Интеграция с внешними LLM API, обработка текстовых запросов
- **Возможности:**
  - Анализ структуры данных
  - Оценка качества данных
  - Генерация рекомендаций
  - Создание DDL скриптов

### Custom LLM (Python)
- **Порт:** 8124
- **Функции:** Обработка запросов с помощью LLM
- **Возможности:**
  - Генерация DDL скриптов
  - Анализ данных
  - Обработка текстовых запросов

### Хранилища данных
- **PostgreSQL** (порт 5432) - Основная база данных
- **MinIO** (порт 9000) - Объектное хранилище файлов
- **Redis** (порт 6379) - Кэширование

## Запуск системы

### Требования
- Docker
- Docker Compose

### Быстрый старт

```bash
# Клонирование репозитория
git clone <repository-url>
cd AIED_baceknd

# Запуск всех сервисов
docker-compose up -d

# Проверка статуса
docker-compose ps

# Тестирование интеграции
./test_integration.sh
```

### Доступные сервисы

После запуска система будет доступна по следующим адресам:

- **Frontend**: http://localhost:3001
- **AI Backend API**: http://localhost:8080
- **LLM Service**: http://localhost:8124
- **MinIO Console**: http://localhost:9001
- **ClickHouse**: http://localhost:8123 (HTTP), http://localhost:9002 (Native)

### Проверка здоровья сервисов

```bash
# Проверка статуса всех контейнеров
docker-compose ps

# Просмотр логов
docker-compose logs -f

# Проверка конкретного сервиса
curl http://localhost:8080/api/v1/health
curl http://localhost:8124/api/v1/health
docker compose ps

# Тестирование системы
./test_llm_integration.sh
```

### Доступ к сервисам

- **Frontend:** http://localhost:3001
- **API Gateway:** http://localhost:8080
- **Data Analysis Service:** http://localhost:8083
- **Custom LLM:** http://localhost:8124
- **MinIO Console:** http://localhost:9001

## Использование

### 1. Загрузка файла
```bash
curl -X POST http://localhost:8080/api/v1/files/upload \
  -F "file=@data.csv" \
  -F "user_id=test_user"
```

### 2. Запуск анализа
```bash
curl -X POST http://localhost:8080/api/v1/analysis/start \
  -H "Content-Type: application/json" \
  -d '{
    "file_id": "file_123",
    "user_id": "test_user",
    "file_path": "bucket/data.csv"
  }'
```

### 3. Проверка статуса
```bash
curl http://localhost:8080/api/v1/analysis/status/{analysis_id}
```

## API Endpoints

### API Gateway

#### POST /api/v1/files/upload
Загрузка файла в систему.

**Параметры:**
- `file` (multipart/form-data) - Файл для загрузки
- `user_id` (query) - ID пользователя

**Ответ:**
```json
{
  "file_id": "file_123",
  "status": "success",
  "message": "Файл загружен успешно",
  "storage_path": "bucket/data.csv"
}
```

#### POST /api/v1/analysis/start
Запуск анализа данных.

**Тело запроса:**
```json
{
  "file_id": "file_123",
  "user_id": "test_user",
  "file_path": "bucket/data.csv"
}
```

**Ответ:**
```json
{
  "analysis_id": "analysis_123",
  "status": "started",
  "message": "Анализ запущен"
}
```

#### GET /api/v1/analysis/status/{analysis_id}
Получение статуса анализа.

**Ответ:**
```json
{
  "analysis_id": "analysis_123",
  "status": "completed",
  "result": {
    "data_quality_score": 0.85,
    "storage_recommendation": {...},
    "table_schema": {...},
    "ddl_metadata": {...}
  }
}
```

### Custom LLM API

#### POST /api/v1/process
Обработка запросов с помощью LLM.

**Тело запроса:**
```json
{
  "user_query": "Проанализируй данные",
  "source_config": {"type": "text"},
  "target_config": {"type": "response"},
  "operation_type": "data_analysis"
}
```

**Ответ:**
```json
{
  "pipeline_id": "pipeline_123",
  "status": "success",
  "message": "Результат анализа",
  "metadata": {
    "processing_time": 0.123,
    "model": "custom-llm-v1.0"
  }
}
```

## Тестирование

### Автоматическое тестирование
```bash
./test_llm_integration.sh
```

### Ручное тестирование

1. **Проверка здоровья сервисов:**
```bash
curl http://localhost:8080/health
curl http://localhost:8083/health
curl http://localhost:8124/health
```

2. **Тестирование LLM:**
```bash
curl -X POST http://localhost:8124/api/v1/process \
  -H "Content-Type: application/json" \
  -d '{
    "user_query": "Привет",
    "source_config": {"type": "text"},
    "target_config": {"type": "response"},
    "operation_type": "text_generation"
  }'
```

## Мониторинг

### Логи сервисов
```bash
# Логи всех сервисов
docker compose logs

# Логи конкретного сервиса
docker compose logs api-gateway
docker compose logs data-analysis-service
docker compose logs custom-llm
```

### Статус контейнеров
```bash
docker compose ps
```

### Использование ресурсов
```bash
docker stats
```

## Устранение неполадок

### Проблемы с запуском
1. Проверьте, что все порты свободны
2. Убедитесь, что Docker запущен
3. Проверьте логи: `docker compose logs`

### Проблемы с LLM
1. Проверьте доступность: `curl http://localhost:8124/health`
2. Проверьте логи: `docker compose logs custom-llm`

### Проблемы с анализом данных
1. Проверьте MinIO: `curl http://localhost:9000/minio/health/live`
2. Проверьте PostgreSQL: `docker exec aien_postgres pg_isready`

## Разработка

### Структура проекта
```
AIED_baceknd/
├── api-geteway/          # API Gateway (Go)
├── data-analysis-service/ # Data Analysis Service (Go)
├── file-service/         # File Service (Go)
├── LLM/                  # Custom LLM (Python)
├── frontend/             # Frontend (React)
├── docker-compose.yml    # Docker Compose конфигурация
├── test_llm_integration.sh # Тесты интеграции
└── README.md             # Этот файл
```

### Добавление новых функций

1. **Новый endpoint в API Gateway:**
   - Добавьте метод в `api-geteway/cmd/http_main.go`
   - Обновите маршруты в `main()`

2. **Новая функциональность в LLM:**
   - Добавьте новый тип операции в `LLM/app.py`
   - Реализуйте обработку в `generate_response()`

3. **Новый анализ в Data Analysis Service:**
   - Добавьте метод в `data-analysis-service/cmd/main.go`
   - Интегрируйте с Custom LLM

## Лицензия

Проект разработан для внутреннего использования.

## Поддержка

Для вопросов и поддержки обращайтесь к команде разработки.