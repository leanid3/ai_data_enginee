FROM python:3.11-slim

WORKDIR /app

# Установка системных зависимостей
RUN apt-get update && apt-get install -y \
    gcc \
    && rm -rf /var/lib/apt/lists/*

# Копирование файлов
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY app.py .

# Создание пользователя для безопасности
RUN useradd -m -u 1000 llmuser && chown -R llmuser:llmuser /app
USER llmuser

# Переменные окружения
ENV FLASK_APP=app.py
ENV FLASK_ENV=production
ENV PORT=8124

EXPOSE 8124

# Команда запуска
CMD ["python", "app.py"]
