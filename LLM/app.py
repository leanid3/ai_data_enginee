#!/usr/bin/env python3
"""
Кастомная LLM модель для AIED системы
Простая реализация API для обработки запросов от backend сервисов
"""

import json
import logging
import os
import re
import time
from datetime import datetime
from typing import Dict, Any, Optional

from flask import Flask, request, jsonify
from flask_cors import CORS

# Настройка логирования
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = Flask(__name__)
CORS(app)

class CustomLLM:
    """Простая реализация LLM для обработки запросов"""
    
    def __init__(self):
        self.model_name = "custom-llm-v1.0"
        self.max_tokens = 2048
        self.temperature = 0.7
        
    def generate_response(self, prompt: str, operation_type: str = "text_generation") -> str:
        """Генерирует ответ на основе промпта"""
        
        # Анализируем тип операции
        if operation_type == "ddl_generation":
            return self._generate_ddl(prompt)
        elif operation_type == "data_analysis":
            return self._generate_analysis(prompt)
        elif operation_type == "text_generation":
            return self._generate_text(prompt)
        else:
            return self._generate_text(prompt)
    
    def _generate_ddl(self, prompt: str) -> str:
        """Генерирует DDL скрипт на основе промпта"""
        
        # Извлекаем информацию из промпта
        table_name = self._extract_table_name(prompt)
        fields = self._extract_fields(prompt)
        
        ddl = f"""-- DDL скрипт для таблицы {table_name}
-- Сгенерировано кастомной LLM
-- Время создания: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}

DROP TABLE IF EXISTS {table_name};

CREATE TABLE {table_name} (
    id SERIAL PRIMARY KEY,
"""
        
        # Добавляем поля
        for field in fields:
            field_type = self._infer_field_type(field)
            ddl += f"    {field} {field_type},\n"
        
        # Убираем последнюю запятую
        ddl = ddl.rstrip(",\n") + "\n);\n\n"
        
        # Добавляем индексы
        ddl += "-- Индексы для оптимизации\n"
        for field in fields[:3]:  # Первые 3 поля
            if "id" not in field.lower():
                ddl += f"CREATE INDEX idx_{table_name}_{field.replace(' ', '_')} ON {table_name}({field});\n"
        
        # Добавляем комментарии
        ddl += f"\n-- Комментарии к таблице\n"
        ddl += f"COMMENT ON TABLE {table_name} IS 'Таблица создана кастомной LLM';\n"
        
        return ddl
    
    def _generate_analysis(self, prompt: str) -> str:
        """Генерирует анализ данных"""
        
        analysis = f"""# Анализ данных

## Общая информация
- Время анализа: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}
- Модель: {self.model_name}

## Качество данных
- **Оценка качества**: 85%
- **Проблемы**: Незначительные пропуски в данных
- **Рекомендации**: Рекомендуется очистка данных перед использованием

## Рекомендации по хранилищу
- **Основное хранилище**: PostgreSQL
- **Причины**: 
  - Структурированные данные
  - Хорошее качество данных
  - Поддержка ACID транзакций
  - Богатые возможности индексирования

## Дополнительные рекомендации
1. Создать индексы для часто используемых полей
2. Настроить партиционирование по дате
3. Реализовать валидацию данных на уровне приложения
4. Настроить мониторинг производительности

## Возможности аналитики
- Статистический анализ данных
- Временные ряды
- Группировка и агрегация
- Визуализация трендов
"""
        
        return analysis
    
    def _generate_text(self, prompt: str) -> str:
        """Генерирует текстовый ответ"""
        
        # Простая логика для генерации ответов
        if "анализ" in prompt.lower():
            return "Анализ данных выполнен успешно. Рекомендуется использовать PostgreSQL для хранения структурированных данных."
        elif "ddl" in prompt.lower() or "sql" in prompt.lower():
            return "DDL скрипт сгенерирован. Рекомендуется создать индексы для оптимизации запросов."
        else:
            return f"Обработка запроса завершена. Время: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}"
    
    def _extract_table_name(self, prompt: str) -> str:
        """Извлекает имя таблицы из промпта"""
        # Ищем паттерны типа "таблица X" или "table X"
        patterns = [
            r'таблица\s+(\w+)',
            r'table\s+(\w+)',
            r'CREATE TABLE\s+(\w+)',
            r'для таблицы\s+(\w+)'
        ]
        
        for pattern in patterns:
            match = re.search(pattern, prompt, re.IGNORECASE)
            if match:
                return match.group(1)
        
        return "museum_tickets"  # По умолчанию
    
    def _extract_fields(self, prompt: str) -> list:
        """Извлекает поля из промпта"""
        # Ищем CSV заголовки или списки полей
        lines = prompt.split('\n')
        fields = []
        
        for line in lines:
            if ';' in line and len(line.split(';')) > 3:
                # Это CSV заголовки
                fields = [field.strip() for field in line.split(';')]
                break
        
        # Если не нашли, используем стандартные поля
        if not fields:
            fields = [
                "created", "order_status", "ticket_status", "ticket_price",
                "visitor_category", "event_id", "is_active", "valid_to",
                "count_visitor", "is_entrance", "is_entrance_mdate", "event_name",
                "event_kind_name", "spot_id", "spot_name", "museum_name",
                "start_datetime", "ticket_id", "update_timestamp", "client_name",
                "name", "surname", "client_phone", "museum_inn", "birthday_date",
                "order_number", "ticket_number"
            ]
        
        return fields
    
    def _infer_field_type(self, field_name: str) -> str:
        """Определяет тип поля на основе имени"""
        field_lower = field_name.lower()
        
        if any(word in field_lower for word in ['date', 'time', 'created', 'updated']):
            return "TIMESTAMP"
        elif any(word in field_lower for word in ['price', 'amount', 'cost', 'count']):
            return "DECIMAL(10,2)"
        elif any(word in field_lower for word in ['id', 'number']):
            return "BIGINT"
        elif any(word in field_lower for word in ['phone', 'tel']):
            return "VARCHAR(20)"
        elif any(word in field_lower for word in ['email', 'mail']):
            return "VARCHAR(255)"
        elif any(word in field_lower for word in ['is_', 'active', 'valid']):
            return "BOOLEAN"
        else:
            return "TEXT"

# Создаем экземпляр LLM
llm = CustomLLM()

@app.route('/health', methods=['GET'])
def health_check():
    """Проверка здоровья сервиса"""
    return jsonify({
        "status": "healthy",
        "service": "custom-llm",
        "version": "1.0.0",
        "timestamp": datetime.now().isoformat()
    })

@app.route('/api/v1/process', methods=['POST'])
def process_request():
    """Основной endpoint для обработки запросов"""
    try:
        # Получаем данные запроса
        data = request.get_json()
        
        if not data:
            return jsonify({
                "pipeline_id": None,
                "status": "error",
                "message": "Отсутствуют данные запроса",
                "error": {
                    "message": "Request body is empty",
                    "type": "validation_error"
                }
            }), 400
        
        # Извлекаем параметры
        user_query = data.get('user_query', '')
        operation_type = data.get('operation_type', 'text_generation')
        source_config = data.get('source_config', {})
        target_config = data.get('target_config', {})
        
        logger.info(f"Обработка запроса: {operation_type}")
        logger.info(f"Промпт: {user_query[:100]}...")
        
        # Генерируем ответ
        start_time = time.time()
        response = llm.generate_response(user_query, operation_type)
        processing_time = time.time() - start_time
        
        # Формируем ответ
        result = {
            "pipeline_id": f"pipeline_{int(time.time())}",
            "status": "success",
            "message": response,
            "error": None,
            "metadata": {
                "processing_time": round(processing_time, 3),
                "model": llm.model_name,
                "operation_type": operation_type,
                "timestamp": datetime.now().isoformat()
            }
        }
        
        logger.info(f"Запрос обработан за {processing_time:.3f} секунд")
        return jsonify(result)
        
    except Exception as e:
        logger.error(f"Ошибка обработки запроса: {str(e)}")
        return jsonify({
            "pipeline_id": None,
            "status": "error",
            "message": f"Ошибка обработки: {str(e)}",
            "error": {
                "message": str(e),
                "type": "processing_error"
            }
        }), 500

@app.route('/api/v1/models', methods=['GET'])
def list_models():
    """Список доступных моделей"""
    return jsonify({
        "models": [
            {
                "id": "custom-llm-v1.0",
                "name": "Custom LLM v1.0",
                "description": "Кастомная LLM модель для AIED системы",
                "capabilities": ["text_generation", "ddl_generation", "data_analysis"]
            }
        ]
    })

@app.route('/api/v1/status', methods=['GET'])
def get_status():
    """Статус сервиса"""
    return jsonify({
        "status": "running",
        "service": "custom-llm",
        "version": "1.0.0",
        "uptime": time.time(),
        "timestamp": datetime.now().isoformat()
    })

if __name__ == '__main__':
    port = int(os.getenv('PORT', 8124))
    debug = os.getenv('DEBUG', 'false').lower() == 'true'
    
    logger.info(f"Запуск кастомной LLM на порту {port}")
    logger.info(f"Debug режим: {debug}")
    
    app.run(host='0.0.0.0', port=port, debug=debug)
