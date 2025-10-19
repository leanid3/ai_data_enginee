# AI Data Engineer — Modular LLM Service (MVP)

## Установка
```bash
python -m venv .venv
. .venv/Scripts/activate  
pip install -U pip
pip install -r requirements.txt

ИЛИ 

conda deactivate && conda create -n hahatonLLM python=3.10 -y && pip install -r requirements.txt 
```

## Конфигурация
Все конфиги в  `.env`
```env
# OpenRouter
OPENROUTER_API_KEY=sk-... #API-ключ
OPENROUTER_BASE_URL=https://openrouter.ai/api/v1
DEFAULT_MODEL=openrouter/auto

# Service metadata (опционально)
SITE_URL=http://localhost:8124
SITE_NAME=AI Data Engineer
```

## Запуск сервиса
```bash
python -m uvicorn main:app --host 0.0.0.0 --port 8124 --reload
```
Сервис поднимется по адресу `http://localhost:8124`.

## Эндпоинты
- POST `/api/v1/process` — основная обработка запроса
- POST `/api/v1/analyze-file` — анализ CSV/JSON файла
- GET `/api/v1/pipeline/{id}` — получить сохранённый пайплайн
- POST `/api/v1/pipeline/{id}/execute` — запуск (заглушка)
- GET `/api/v1/health` — статус
- GET `/api/v1/agents` — список агентов

Документация Swagger: `http://localhost:8124/docs`

## Быстрый тест (ручной)
1) Проверить здоровье:
```bash
curl http://localhost:8124/api/v1/health
```
2) Запустить процесс:
```bash
curl -X POST http://localhost:8124/api/v1/process \
  -H "Content-Type: application/json" \
  -d '{
    "user_query": "Проанализировать данные о продажах и построить ETL",
    "source_config": {"type": "csv", "file_path": "sales.csv"},
    "target_config": {"type": "clickhouse"},
    "operation_type": "full_process"
  }'
```
3) Анализ файла:
```bash
curl -X POST "http://localhost:8124/api/v1/analyze-file" \
  -H "accept: application/json" \
  -H "Content-Type: multipart/form-data" \
  -F "file=@path/to/data.csv"
```
Ответы включают `pipeline_id`, который можно запросить через GET `/api/v1/pipeline/{id}`.

## Тестирование 

Пример фиктивного клиента для тестов:
```python
class FakeLLM:
    async def chat(self, messages, *, model, temperature, max_tokens, response_format):
        return {"choices": [{"message": {"content": "{\\"ok\\": true}"}}]}
```
Поменяйте `app.state.llm_client = FakeLLM()` 

## Структура
```
core/
  agents.py          # Роли, промпты, базовые агенты 
  llm_client.py      # OpenRouterClient 
  router.py          # SmartRouter 
api/
  main.py            # FastAPI и эндпоинты
```

## Тестирование агентов
Для проверки работы всех агентов системы используйте тестовый скрипт:

```bash
python test_agents.py
```

Скрипт тестирует все 7 агентов:
- **ROUTER** — маршрутизация запросов
- **DATA_ANALYZER** — анализ структуры данных
- **DB_SELECTOR** — выбор системы хранения
- **DDL_GENERATOR** — генерация DDL схем
- **ETL_BUILDER** — построение ETL пайплайнов
- **QUERY_OPTIMIZER** — оптимизация SQL запросов
- **REPORT_GENERATOR** — создание отчетов



Результаты сохраняются в `test_results.json`.

## Полезные команды
- Запуск  приложения: `python -m uvicorn main:app --host 0.0.0.0 --port 8124 --reload`
- Тестирование агентов: `python test_agents.py`



