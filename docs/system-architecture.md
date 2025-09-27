# Диаграмма архитектуры AIEN Backend

## Общая архитектура системы

```mermaid
graph TB
    subgraph "Frontend Layer"
        FE[Frontend Application]
    end
    
    subgraph "API Gateway Layer"
        AG[API Gateway<br/>:8080 HTTP<br/>:50051 gRPC]
    end
    
    subgraph "Business Logic Layer"
        DP[Data Profiler<br/>:50052 gRPC]
        FS[File Service]
        PS[Profile Service]
    end
    
    subgraph "External Services"
        LLM[Ollama LLM<br/>:11434 HTTP]
    end
    
    subgraph "Data Sources"
        CSV[CSV Files]
        JSON[JSON Files]
        XML[XML Files]
        DB[(Database)]
        STREAM[Stream Data]
    end
    
    subgraph "Target Systems"
        PG[(PostgreSQL)]
        CH[(ClickHouse)]
        HDFS[HDFS]
        ES[(Elasticsearch)]
        S3[S3 Storage]
    end
    
    FE -->|HTTP Requests| AG
    AG -->|gRPC Calls| DP
    DP -->|HTTP Requests| LLM
    AG -->|HTTP Requests| LLM
    
    CSV -->|Upload| AG
    JSON -->|Upload| AG
    XML -->|Upload| AG
    DB -->|Connect| AG
    STREAM -->|Connect| AG
    
    AG -->|DDL Scripts| PG
    AG -->|DDL Scripts| CH
    AG -->|DDL Scripts| HDFS
    AG -->|DDL Scripts| ES
    AG -->|DDL Scripts| S3
```

## Поток данных

```mermaid
sequenceDiagram
    participant F as Frontend
    participant AG as API Gateway
    participant DP as Data Profiler
    participant LLM as Ollama LLM
    participant T as Target System
    
    F->>AG: 1. Upload File (CSV/JSON/XML)
    AG->>DP: 2. Process File
    DP->>DP: 3. Analyze Data
    DP-->>AG: 4. Data Profile
    
    F->>AG: 5. Create Pipeline
    AG->>LLM: 6. Generate DDL
    LLM-->>AG: 7. DDL Script
    AG-->>F: 8. Pipeline Created
    
    F->>AG: 9. Execute Pipeline
    AG->>T: 10. Deploy DDL
    T-->>AG: 11. Data Loaded
    AG-->>F: 12. Pipeline Complete
```

## Компоненты системы

### API Gateway
- **Порт**: 8080 (HTTP), 50051 (gRPC)
- **Функции**:
  - Единственная точка входа
  - Проксирование запросов
  - Интеграция с LLM
  - Управление пайплайнами

### Data Profiler
- **Порт**: 50052 (gRPC)
- **Функции**:
  - Загрузка и обработка файлов
  - Анализ структуры данных
  - Профилирование качества данных
  - Преобразование в единый формат

### Ollama LLM
- **Порт**: 11434 (HTTP)
- **Функции**:
  - Анализ профиля данных
  - Генерация DDL скриптов
  - Рекомендации по оптимизации
  - Создание индексов и ограничений

## Поддерживаемые форматы данных

| Формат | Размер | Описание | Использование |
|--------|--------|-----------|---------------|
| CSV | до 1GB | Структурированные данные | Транзакционные системы |
| JSON | до 1GB | Полуструктурированные данные | API и веб-приложения |
| XML | до 1GB | Иерархические данные | Документооборот |
| Database | - | Прямое подключение | Существующие БД |
| Stream | - | Потоковые данные | Real-time обработка |

## Целевые системы хранения

| Система | Тип | Оптимизация | Использование |
|---------|-----|-------------|----------------|
| PostgreSQL | OLTP | ACID, индексы | Транзакционные данные |
| ClickHouse | OLAP | Колоночное хранение | Аналитика и отчеты |
| HDFS | Big Data | Распределенное хранение | Большие объемы |
| Elasticsearch | Search | Поисковые индексы | Логи и поиск |
| S3 | Object Storage | Объектное хранение | Архивирование |

## Безопасность и масштабируемость

### Безопасность
- Аутентификация через API ключи
- Шифрование данных в транзите
- Изоляция сервисов в Docker
- Валидация входных данных

### Масштабируемость
- Горизонтальное масштабирование сервисов
- Асинхронная обработка больших файлов
- Кэширование результатов анализа
- Оптимизация gRPC соединений