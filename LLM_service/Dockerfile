FROM python:3.10-slim

# Установка системных зависимостей
RUN apt-get update && apt-get install -y \
    gcc \
    g++ \
    && rm -rf /var/lib/apt/lists/*

# Установка рабочей директории
WORKDIR /app

# Копирование файлов зависимостей
COPY requirements.txt requirements.txt

# Установка Python зависимостей
RUN pip install --no-cache-dir -r requirements.txt

# Копирование исходного кода
COPY . .

# Создание директории для логов
RUN mkdir -p logs

# Открытие порта
EXPOSE 8124

# Переменные окружения
ENV PYTHONPATH=/app
ENV PORT=8124

# Команда запуска
CMD ["python", "-m", "uvicorn", "API.app:app", "--host", "0.0.0.0", "--port", "8124"]
