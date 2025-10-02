# Настройка кастомной LLM

## Обзор изменений

Проект был успешно обновлен для использования вашей кастомной LLM вместо Ollama. Все сервисы теперь настроены для работы с вашим API gateway.

## Изменения в коде

### 1. API Gateway (`api-geteway/`)
- **config/config.go**: Обновлена конфигурация для кастомной LLM
- **internal/clients/llm_client.go**: Переписан клиент для работы с вашим API

### 2. LLM Service (`llm-service/`)
- **config/config.go**: Обновлена конфигурация
- **internal/custom_llm/custom_llm_client.go**: Новый клиент для кастомной LLM
- **cmd/main.go**: Обновлен для использования нового клиента

### 3. Data Analysis Service (`data-analysis-service/`)
- **cmd/main.go**: Полностью переписан для работы с кастомной LLM
- **internal/custom_llm/custom_llm_client.go**: Новый клиент

### 4. Airflow DAG (`airflow/dags/data_analysis_dag.py`)
- Обновлен для использования кастомной LLM

### 5. Docker Compose (`docker-compose.yml`)
- Удален сервис Ollama
- Добавлены переменные окружения для кастомной LLM

## Переменные окружения

### Для всех сервисов:
```bash
CUSTOM_LLM_URL=http://localhost:8124/api/v1/process
CUSTOM_LLM_API_KEY=your_api_key_here
LLM_MODEL=openrouter/auto
```

### Для API Gateway:
```bash
LLM_BASE_URL=http://localhost:8124/api/v1/process
LLM_API_KEY=your_api_key_here
LLM_MODEL=openrouter/auto
```

## Формат запросов к кастомной LLM

Все сервисы теперь отправляют запросы в следующем формате:

```json
{
  "user_query": "Ваш запрос к LLM",
  "source_config": {
    "type": "text"
  },
  "target_config": {
    "type": "response"
  },
  "operation_type": "text_generation"
}
```

## Ожидаемый формат ответа

Кастомная LLM должна возвращать ответы в формате:

```json
{
  "pipeline_id": "optional_id",
  "status": "success",
  "message": "Ответ от LLM",
  "error": null
}
```

## Запуск системы

1. Убедитесь, что ваша кастомная LLM запущена на `http://localhost:8124`
2. Запустите Docker Compose:
   ```bash
   docker compose up -d
   ```

### После изменений в коде

Если вы внесли изменения в код сервисов, пересоберите образы:

```bash
# Пересборка всех сервисов
docker compose build

# Пересборка конкретного сервиса
docker compose build llm-service
docker compose build data-analysis-service

# Перезапуск сервисов
docker compose up -d
```

### Важные изменения в docker-compose.yml

- ✅ Удален сервис Ollama
- ✅ Обновлены переменные окружения для кастомной LLM
- ✅ Исправлены отступы YAML
- ✅ Добавлен frontend сервис
- ✅ Исправлена структура зависимостей

## Тестирование

Для тестирования интеграции используйте:

```bash
# Проверка здоровья API Gateway
curl http://localhost:8080/health

# Тест анализа данных
curl -X POST http://localhost:8083/analyze \
  -H "Content-Type: application/json" \
  -d '{"file_id": "test", "user_id": "user1", "file_path": "test.csv"}'
```

## Важные замечания

1. **Таймауты**: Установлены увеличенные таймауты (10 минут) для обработки больших файлов
2. **Fallback**: Все сервисы имеют fallback механизмы при недоступности LLM
3. **Логирование**: Добавлено подробное логирование для отладки
4. **Безопасность**: API ключи передаются через переменные окружения

## Устранение неполадок

1. **LLM недоступна**: Проверьте, что ваша кастомная LLM запущена и доступна
2. **Ошибки аутентификации**: Убедитесь, что API ключ установлен правильно
3. **Таймауты**: Увеличьте таймауты в конфигурации при необходимости
