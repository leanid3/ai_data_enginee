# AI Data Engineer Backend

Backend сервис для AI Data Engineer MVP, построенный на Go с использованием Clean Architecture.

## Архитектура

Проект следует принципам Clean Architecture с четким разделением слоев:

```
ai-data-engineer-backend/
├── cmd/server/           # Точка входа приложения
├── internal/
│   ├── api/             # HTTP слой
│   │   ├── handlers/    # HTTP обработчики
│   │   └── middleware/  # Middleware
│   ├── service/         # Бизнес-логика
│   ├── models/          # Модели данных
│   ├── repository/      # Слой доступа к данным
│   └── config/          # Конфигурация
├── pkg/                 # Переиспользуемые пакеты
│   ├── logger/         # Структурированное логирование
│   ├── validator/      # Валидация
│   └── utils/          # Утилиты
├── configs/            # Конфигурационные файлы
└── scripts/            # Скрипты и шаблоны
```

## Возможности

- **Анализ файлов**: Поддержка CSV, JSON, XML форматов
- **LLM интеграция**: Анализ данных с помощью языковых моделей
- **Генерация DDL**: Автоматическое создание схем БД
- **ETL пайплайны**: Построение и выполнение пайплайнов данных
- **Airflow интеграция**: Генерация DAG для Apache Airflow
- **Множественные БД**: PostgreSQL, ClickHouse, HDFS
- **Микросервисная архитектура**: REST API с четким разделением ответственности

## API Endpoints

### Файлы
- `POST /api/v1/files/upload` - Загрузка и анализ файла
- `GET /api/v1/files/:id` - Получение информации о файле
- `DELETE /api/v1/files/:id` - Удаление файла
- `GET /api/v1/files` - Список файлов пользователя

### Анализ данных
- `POST /api/v1/analysis/start` - Запуск анализа
- `GET /api/v1/analysis/:id/status` - Статус анализа
- `GET /api/v1/analysis/:id/result` - Результат анализа
- `GET /api/v1/analysis` - Список анализов

### Пайплайны
- `POST /api/v1/pipelines` - Создание пайплайна
- `GET /api/v1/pipelines/:id` - Получение пайплайна
- `POST /api/v1/pipelines/:id/execute` - Выполнение пайплайна
- `DELETE /api/v1/pipelines/:id` - Удаление пайплайна
- `GET /api/v1/pipelines` - Список пайплайнов

### Health Check
- `GET /api/v1/health` - Проверка состояния сервиса
- `POST /api/v1/databases/test` - Тестирование подключения к БД

## Быстрый старт

### Предварительные требования

- Go 1.21+
- Docker и Docker Compose
- Make (опционально)

### Установка

1. **Клонируйте репозиторий**:
```bash
git clone <repository-url>
cd ai-data-engineer-backend
```

2. **Установите зависимости**:
```bash
make deps
# или
go mod download
```

3. **Настройте конфигурацию**:
```bash
cp .env.example .env
# Отредактируйте .env файл при необходимости
```

### Запуск

#### Локальная разработка

```bash
# Запуск всех сервисов
make docker-up

# Запуск только backend
make run

# Запуск с hot reload (требует air)
make dev
```

#### Production

```bash
# Сборка и запуск
make build
make run

# Или через Docker
make docker-build
make docker-run
```

### Доступные команды

```bash
make help                    # Показать все команды
make build                   # Собрать приложение
make run                     # Запустить локально
make test                    # Запустить тесты
make docker-up               # Запустить все сервисы
make docker-down             # Остановить все сервисы
make docker-logs              # Показать логи
make lint                    # Проверить код
make fmt                     # Форматировать код
make clean                   # Очистить артефакты
```

## Конфигурация

### Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `SERVER_PORT` | Порт сервера | `8080` |
| `CONFIG_PATH` | Путь к конфигу | `configs/config.yaml` |
| `POSTGRES_HOST` | Хост PostgreSQL | `postgres` |
| `POSTGRES_PORT` | Порт PostgreSQL | `5432` |
| `CLICKHOUSE_HOST` | Хост ClickHouse | `clickhouse` |
| `LLM_BASE_URL` | URL LLM сервиса | `http://custom-llm:8124/api/v1/process` |
| `LOG_LEVEL` | Уровень логирования | `info` |

### Конфигурационные файлы

- `configs/config.yaml` - Локальная конфигурация
- `configs/config.prod.yaml` - Production конфигурация

## Разработка

### Структура проекта

```
internal/
├── api/                    # HTTP слой
│   ├── handlers/          # Обработчики запросов
│   │   ├── file_handler.go
│   │   ├── pipeline_handler.go
│   │   ├── analyze_handler.go
│   │   └── health_handler.go
│   ├── middleware/        # Middleware
│   │   ├── cors.go
│   │   ├── logger.go
│   │   └── error_handler.go
│   └── routes.go          # Маршруты
├── service/               # Бизнес-логика
│   ├── interfaces.go     # Интерфейсы сервисов
│   ├── requests.go       # Request/Response модели
│   └── ...               # Реализации сервисов
├── models/                # Модели данных
│   ├── request.go
│   ├── response.go
│   ├── pipeline.go
│   ├── schema.go
│   └── errors.go
├── repository/            # Слой доступа к данным
│   ├── interface.go      # Интерфейсы репозиториев
│   ├── postgres.go       # PostgreSQL репозиторий
│   ├── clickhouse.go     # ClickHouse репозиторий
│   └── pipeline_storage.go
└── config/               # Конфигурация
    └── config.go
```

### Принципы разработки

1. **Clean Architecture**: Четкое разделение слоев
2. **Dependency Injection**: Внедрение зависимостей через интерфейсы
3. **Error Handling**: Структурированная обработка ошибок
4. **Logging**: Структурированное логирование с контекстом
5. **Testing**: Unit и integration тесты
6. **Documentation**: Документированный код

### Добавление новых функций

1. **Создайте модели** в `internal/models/`
2. **Определите интерфейсы** в `internal/service/`
3. **Реализуйте сервисы** в `internal/service/`
4. **Создайте handlers** в `internal/api/handlers/`
5. **Добавьте маршруты** в `internal/api/routes.go`
6. **Напишите тесты** для всех компонентов

## Мониторинг и логирование

### Логирование

- **Формат**: JSON (production) / Pretty (development)
- **Уровни**: DEBUG, INFO, WARN, ERROR, FATAL
- **Контекст**: Request ID, User ID, Service
- **Структурированные поля**: timestamp, level, message, fields

### Health Checks

- **Endpoint**: `/api/v1/health`
- **Проверки**: Сервис, БД, LLM, Storage
- **Статусы**: healthy, unhealthy, degraded

### Метрики

- HTTP запросы (количество, время ответа, ошибки)
- БД операции (подключения, запросы, ошибки)
- LLM запросы (количество, время ответа, токены)
- Файловые операции (загрузки, размеры, типы)

## Развертывание

### Docker

```bash
# Сборка образа
docker build -t ai-data-engineer-backend .

# Запуск контейнера
docker run -p 8080:8080 \
  -e POSTGRES_HOST=postgres \
  -e LLM_BASE_URL=http://llm:8124/api/v1/process \
  ai-data-engineer-backend
```

### Docker Compose

```bash
# Запуск всех сервисов
docker-compose up -d

# Просмотр логов
docker-compose logs -f ai-backend

# Остановка
docker-compose down
```

### Kubernetes

```yaml
# Пример deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ai-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: ai-backend
  template:
    metadata:
      labels:
        app: ai-backend
    spec:
      containers:
      - name: ai-backend
        image: ai-data-engineer-backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: POSTGRES_HOST
          value: "postgres-service"
        - name: LLM_BASE_URL
          value: "http://llm-service:8124/api/v1/process"
```

## Troubleshooting

### Частые проблемы

1. **Ошибка подключения к БД**:
   - Проверьте переменные окружения
   - Убедитесь, что БД запущена
   - Проверьте сетевые настройки

2. **LLM сервис недоступен**:
   - Проверьте URL и API ключ
   - Убедитесь, что сервис запущен
   - Проверьте таймауты

3. **Ошибки валидации**:
   - Проверьте формат входных данных
   - Убедитесь в корректности типов
   - Проверьте обязательные поля

### Логи и отладка

```bash
# Просмотр логов
docker-compose logs -f ai-backend

# Логи с фильтрацией
docker-compose logs -f ai-backend | grep ERROR

# Отладка подключений
docker-compose exec ai-backend wget -O- http://postgres:5432
```

## Лицензия

MIT License

## Контрибьюция

1. Fork репозиторий
2. Создайте feature branch
3. Сделайте commit изменений
4. Push в branch
5. Создайте Pull Request

## Поддержка

- **Issues**: GitHub Issues
- **Документация**: README.md
- **Примеры**: `/examples` директория

